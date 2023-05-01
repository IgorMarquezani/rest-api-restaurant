package products_lists

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var (
		room models.Room
		pl   models.ProductList
	)

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewDecoder(r.Body).Decode(&pl)

	if pl.Room <= 0 {
		room = models.RoomByItsId(session.ActiveRoom)
		pl.Room = room.Id
	} else {
		room = models.RoomByItsId(pl.Room)
	}

	room.FindGuests()

	if !room.IsOwner(user) {
		if !room.IsGuest(user) {
			http.Error(w, "You are not a guest in that room", http.StatusForbidden)
			return
		}

		if room.GuestPermission(user) < 2 {
			http.Error(w, "You do not have permission", http.StatusForbidden)
			return
		}
	}

	if err := models.InsertProductList(pl); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "Name already in use", http.StatusAlreadyReported)
			return
		}

		if models.ErrPlEmptyName == err.Error() {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.Encode(pl)
}
