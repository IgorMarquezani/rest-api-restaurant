package tabs

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
)

type Register struct {
	user    models.User
	session models.UserSession
	tab     models.Tab
	room    models.Room
}

func (re Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var howManyIsertions uint
	var err error

  // if flusher, ok := w.(http.Flusher); ok {
  //   flusher.Flush()
  // }

	err, re.user, re.session = controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

  err = json.NewDecoder(r.Body).Decode(&re.tab)
  if err != nil {
    panic(err)
  }
	re.room = models.RoomByItsId(re.tab.RoomId)

	if re.room.Id == 0 {
		re.room = models.RoomByItsId(re.session.ActiveRoom)
		re.tab.RoomId = re.room.Id
	}

	if re.tab.PayValue != 0 {
		re.tab.PayValue = 0
	}

	controllers.AllowCrossOrigin(&w, "*")

	if !re.room.IsOwner(re.user) {
		if !re.room.IsGuest(re.user) {
			http.Error(w, "Not a guest in that room", http.StatusUnauthorized)
			return
		}
	}

	if err := models.InsertTab(&re.tab); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "Tab already exists", http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(re.tab.Requests); howManyIsertions++ {
		if re.tab.Requests[i].Quantity <= 0 {
			howManyIsertions--
			continue
		}

		err := models.InsertRequest(re.tab, re.tab.Requests[i])
		if err != nil && database.IsDuplicateKeyError(err.Error()) {
			howManyIsertions--
			continue
		}

		if err != nil {
			models.DeleteTab(re.tab)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if howManyIsertions == 0 {
		models.DeleteTab(re.tab)
		http.Error(w, "Tab not registered because all requests have 0 zero quantity", http.StatusExpectationFailed)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(re.tab)
}
