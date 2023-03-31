package session

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/models"
)

func MustUpdateActiveRoom(w http.ResponseWriter, r *http.Request) {
  var room models.Room

  err, user, session := controllers.VerifySession(r)
  if err != nil {
    panic(err)
  }

  json.NewDecoder(r.Body).Decode(&room)

  if !room.IsOwner(user) {
    if !room.IsGuest(user) {
      http.Error(w, "Not a Guest in that room", http.StatusUnauthorized)
      return
    }
  }

  if err := models.UpdateActiveRoom(&session, room); err != nil {
    if err.Error() == models.NoSessionErr {
      http.Error(w, "Please log in", http.StatusBadRequest)
      return
    }

    http.Error(w, "Iternal server error", http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusAccepted)
}
