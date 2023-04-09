package products

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
	"github.com/gorilla/mux"
)

func GetProduct(w http.ResponseWriter, r *http.Request) {
  var room models.Room

  err, user, session := controllers.VerifySession(r)
  if err != nil {
    http.Error(w, "Not loged in", http.StatusUnauthorized)
    return
  }

  name := mux.Vars(r)["name"]
  if name == "" {
    http.Error(w, "Wrong query parameter", http.StatusBadRequest)
    return
  }

  idStr := mux.Vars(r)["room"]
  id, err := strconv.Atoi(idStr)
  if err != nil {
    http.Error(w, "Invalid room id", http.StatusBadRequest)
    return
  }

  if id > 0 {
    room = models.RoomByItsId(id)
  } else {
    room = models.RoomByItsId(session.ActiveRoom)
  }

  if !room.IsOwner(user) {
    if !room.IsGuest(user) {
      http.Error(w, "Not a guest", http.StatusUnauthorized)
      return
    }
  }

  p, err := models.SelectOneProduct(room.Id, name)
  if err != nil {
    if err.Error() == database.ErrNotfound {
      http.Error(w, err.Error(), http.StatusNoContent)
      return
    }

    http.Error(w, "Iternal server error", http.StatusInternalServerError)
    log.Println(err.Error())
    return
  }

  w.WriteHeader(http.StatusOK)
  encoder := json.NewEncoder(w)
  encoder.SetIndent("", "    ")
  encoder.Encode(p)
}
