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
		product models.Product
		room    models.Room
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

	json.Unmarshal(body, &product)
	if product.ListRoom <= 0 {
		room = models.RoomByItsId(session.ActiveRoom)
		product.ListRoom = room.Id
	} else {
		room = models.RoomByItsId(product.ListRoom)
	}

	if !room.IsOwner(user) {
		if !room.IsGuest(user) {
			http.Error(w, models.ErrNotAGuest, http.StatusBadRequest)
			return
		}

		if room.GuestPermission(user) < 2 {
			http.Error(w, models.ErrInvalidPermission, http.StatusUnauthorized)
			return
		}
	}

	if err := models.InsertProduct(product); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, models.ErrProductNameAlreadyUsed, http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if _, ok := models.RoomProducts[room.Id]; !ok {
		models.RoomProducts[room.Id] = make(models.Products)
	}

	models.RoomProducts[room.Id][product.Name] = product

	w.WriteHeader(http.StatusCreated)
	controllers.EncodeJSON(w, product)
}
