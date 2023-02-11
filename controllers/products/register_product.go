package products

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
)

type RegisterProduct struct {
	product     models.Product
	productList models.ProductList
	jsonRoom    models.Room
}

func (rp RegisterProduct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, user, _ := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user.Room = models.RoomByItsOwner(user.Id)

	json.NewDecoder(r.Body).Decode(&rp.product)

	rp.jsonRoom = models.RoomByItsId(rp.product.ListRoom)
	rp.jsonRoom.FindGuests()

	if !rp.jsonRoom.IsOwner(user) {
		if !rp.jsonRoom.IsGuest(user) {
			http.Error(w, "You are not a guest", http.StatusBadRequest)
			return
		}

		if rp.jsonRoom.GuestPermission(user) < 2 {
			http.Error(w, "You do not have permission for this operation", http.StatusUnauthorized)
			return
		}
	}

	rp.productList.Name = rp.product.ListName
	rp.productList.Room = rp.product.ListRoom

	if err := models.InsertProduct(rp.product, rp.productList); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "Product name already in use", http.StatusConflict)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rp.product)
}
