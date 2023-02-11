package products

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
)

type UpdateProduct struct {
	bothProducts models.UpdatingProduct
	productList  models.ProductList
	jsonRoom     models.Room
}

func (up UpdateProduct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err, user = controllers.VerifySessionCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Room = models.RoomByItsOwner(user.Id)

	json.NewDecoder(r.Body).Decode(&up.bothProducts)

	if err := up.bothProducts.IsIncompatible(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	up.productList.Name = up.bothProducts.Old.ListName
	up.productList.Room = up.bothProducts.Old.ListRoom

	if !up.productList.Exists() {
		http.Error(w, "No such product list", http.StatusBadRequest)
		return
	}

	if !up.productList.IsOnProductList(up.bothProducts.Old) {
		http.Error(w, "No such product to be update", http.StatusBadRequest)
		return
	}

	up.jsonRoom = models.RoomByItsId(up.bothProducts.Old.ListRoom)

	if !up.jsonRoom.IsOwner(user) {
		if !up.jsonRoom.IsGuest(user) {
			http.Error(w, "You are not a guest in that room", http.StatusUnauthorized)
			return
		}

		if up.jsonRoom.GuestPermission(user) < 2 {
			http.Error(w, "You do not have permission for this operation", http.StatusUnauthorized)
			return
		}
	}

	if !up.bothProducts.Old.Exists() {
		http.Error(w, "No such product to be updated", http.StatusNotFound)
		return
	}

	if err := models.UpdateProduct(up.bothProducts, up.productList); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "Name already in use", http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(up.bothProducts.New)
}
