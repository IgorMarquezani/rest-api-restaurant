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

	if err := models.InsertTab(&handler.tab); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "Tab already exists", http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(handler.tab.Requests); i++ {
		err := models.InsertRequest(handler.tab, handler.tab.Requests[i])
		if err != nil {
			models.DeleteTab(handler.tab)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(handler.tab)
}
