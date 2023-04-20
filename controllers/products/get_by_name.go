package products

import (
	"log"
	"net/http"
	"strconv"

	"github.com/api/controllers"
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
		http.Error(w, "Invalid name parameter", http.StatusBadRequest)
		return
	}

	idStr := mux.Vars(r)["room"]
	id, err := strconv.Atoi(idStr)
	if err != nil && idStr != "" {
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
			http.Error(w, models.ErrNotAGuest, http.StatusForbidden)
			return
		}

    if room.GuestPermission(user) < 2 {
			http.Error(w, models.ErrInvalidPermission, http.StatusForbidden)
			return
    }
	}

	product, err := models.SelectOneProduct(room.Id, name)
	if err != nil {
		if err.Error() == models.ErrNoSuchProduct {
			http.Error(w, models.ErrNoSuchProduct, http.StatusNotFound)
			return
		}

		http.Error(w, "Iternal server error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	controllers.EncodeJSON(w, product)
}
