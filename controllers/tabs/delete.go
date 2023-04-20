package tabs

import (
	"log"
	"net/http"
	"strconv"

	"github.com/api/controllers"
	"github.com/api/models"
	"github.com/gorilla/mux"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	controllers.AllowCrossOrigin(&w, r.Header.Get("Origin"))
  w.Header().Set("Acess-Control-Allow-Credentials", "true")
  log.Println(w.Header())

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	numberStr := mux.Vars(r)["number"]
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		http.Error(w, "Invalid number for tab", http.StatusBadRequest)
		return
	}

	roomIdStr := mux.Vars(r)["room"]
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		roomId = models.RoomByItsId(session.ActiveRoom).Id
	}

  room := models.RoomByItsId(roomId)

  if !room.IsOwner(user) {
    if !room.IsGuest(user) {
			http.Error(w, "Not a guest in that room", http.StatusUnauthorized)
			return
    }
  }

	if err := models.DeleteTab(models.Tab{Number: number, RoomId: roomId}); err != nil {
		log.Println(err)
		http.Error(w, "Tab does not exist in this room", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
