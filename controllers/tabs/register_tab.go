package tabs

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var (
		howManyIsertions uint
		tab              models.Tab
		room             models.Room
	)

	// if flusher, ok := w.(http.Flusher); ok {
	//   flusher.Flush()
	// }

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	body, err := controllers.ValidJSONFormat(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Unmarshal(body, &tab)

	room = models.RoomByItsId(tab.RoomId)

	if room.Id == 0 {
		room = models.RoomByItsId(session.ActiveRoom)
		tab.RoomId = room.Id
	}

	if tab.PayValue != 0 {
		tab.PayValue = 0
	}

	controllers.AllowCrossOrigin(&w, "*")

	if !room.IsOwner(user) {
		if !room.IsGuest(user) {
			http.Error(w, "Not a guest in that room", http.StatusUnauthorized)
			return
		}
	}

	if err := models.InsertTab(&tab); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "Tab already exists", http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(tab.Requests); i++ {
		if tab.Requests[i].Quantity <= 0 {
			continue
		}

		err := models.InsertRequest(tab, tab.Requests[i])
		if err != nil && database.IsDuplicateKeyError(err.Error()) {
			request := models.SelectRequest(tab.Requests[i].ProductName, tab.Number, tab.RoomId)
			models.UpdateRequestQuantity(request, uint(request.Quantity+tab.Requests[i].Quantity))
			howManyIsertions++
			continue
		}

		if err != nil {
			models.DeleteTab(tab)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		howManyIsertions++
	}

	if howManyIsertions == 0 {
		models.DeleteTab(tab)
		http.Error(w, "Tab not registered because all requests have 0 zero quantity", http.StatusExpectationFailed)
		return
	}

  sendTabInChan(room.Id, tab)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tab)
}
