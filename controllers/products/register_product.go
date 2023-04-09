package products

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var (
		p    models.Product
		room models.Room
	)

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := controllers.ValidJSONFormat(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	json.Unmarshal(body, &p)
	if p.ListRoom <= 0 {
		room = models.RoomByItsId(session.ActiveRoom)
		p.ListRoom = room.Id
	} else {
		room = models.RoomByItsId(p.ListRoom)
		p.ListRoom = room.Id
	}

	if !room.IsOwner(user) {
		if !room.IsGuest(user) {
			http.Error(w, "You are not a guest", http.StatusBadRequest)
			return
		}

		if room.GuestPermission(user) < 2 {
			http.Error(w, "You do not have permission for this operation", http.StatusUnauthorized)
			return
		}
	}

	if err := models.InsertProduct(p); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "Product name already in use", http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if _, ok := models.RoomProducts[room.Id]; !ok {
		models.RoomProducts[room.Id] = make(models.Products)
	}

	models.RoomProducts[room.Id][p.Name] = p

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.Encode(p)
}
