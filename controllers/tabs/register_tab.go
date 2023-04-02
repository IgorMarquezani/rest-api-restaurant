package tabs

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
)

type HandleTabRegister struct {
	tab      models.Tab
	jsonRoom models.Room
}

func (handler HandleTabRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var howManyIsertions uint

	controllers.AllowCrossOrigin(&w, "*")

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewDecoder(r.Body).Decode(&handler.tab)
	handler.jsonRoom = models.RoomByItsId(handler.tab.RoomId)

	if handler.jsonRoom.Id == 0 {
		handler.jsonRoom = models.RoomByItsId(session.ActiveRoom)
		handler.tab.RoomId = handler.jsonRoom.Id
	}

	if !handler.jsonRoom.IsOwner(user) {
		if !handler.jsonRoom.IsGuest(user) {
			http.Error(w, "Not a guest in that room", http.StatusUnauthorized)
			return
		}
	}

	handler.tab.CalculateValue()
	if err := models.InsertTab(&handler.tab); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "Tab already exists", http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}

	for i := range handler.tab.Requests {
		if handler.tab.Requests[i].Quantity <= 0 {
			continue
		}

		err := models.InsertRequest(handler.tab, handler.tab.Requests[i])
		if err != nil && database.IsDuplicateKeyError(err.Error()) {
			continue
		}

		if err != nil {
			models.DeleteTab(handler.tab)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		howManyIsertions++
	}

	if howManyIsertions == 0 {
		models.DeleteTab(handler.tab)
		http.Error(w, "Tab no registered because as requests have 0 zero quantity", http.StatusExpectationFailed)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(handler.tab)
}
