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

	controllers.AllowCrossOrigin(&w, r.Header.Get("Origin"))
	w.Header().Set("Acess-Control-Allow-Credentials", "true")

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := controllers.ValidJSONFormat(r.Body)
	if err != nil {
		http.Error(w, "JSON format error: "+err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()

	json.Unmarshal(body, &tab)

	room = models.RoomByItsId(session.ActiveRoom)

	tab.RoomId = room.Id
	tab.PayValue = 0

	if ok, _ := room.IsOwnerOrGuest(user); !ok {
		http.Error(w, "Not a guest in that room", http.StatusUnauthorized)
		return
	}

	if err := models.InsertTab(&tab); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "Tab already exists", http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}

	tab.GroupRequests()

	for i := 0; i < len(tab.Requests); i++ {
		if tab.Requests[i].Quantity <= 0 {
			continue
		}

		err := models.InsertRequest(tab, tab.Requests[i])
		if err != nil && database.IsDuplicateKeyError(err.Error()) {
			request, err := models.SelectRequest(tab.Requests[i].ProductName, tab.Number, tab.RoomId)
			if err != nil {
				panic(err)
			}

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

	SendTab(room.Id, tab)

	w.WriteHeader(http.StatusCreated)
	controllers.EncodeJSON(w, tab)
}
