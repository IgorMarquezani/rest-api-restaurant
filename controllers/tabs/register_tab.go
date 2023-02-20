package tabs

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
)

type TabRegister struct {
  tab models.Tab
  jsonRoom models.Room
}

func (tr TabRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, user, _ := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

  json.NewDecoder(r.Body).Decode(&tr.tab)
  tr.jsonRoom = models.RoomByItsId(tr.tab.RoomId)

  if !tr.jsonRoom.IsOwner(user) {
    if !tr.jsonRoom.IsGuest(user) {
      http.Error(w, "Not a guest in that room", http.StatusUnauthorized)
      return
    }
  }

  if err := models.InsertTab(&tr.tab); err != nil {
    if database.IsDuplicateKeyError(err.Error()) {
      http.Error(w, "Tab already exists", http.StatusAlreadyReported)
      return
    }

    http.Error(w, "Internal Server error", http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(tr.tab)
}
