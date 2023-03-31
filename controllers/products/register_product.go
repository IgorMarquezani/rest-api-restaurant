package products

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
)

type HandleProductRegister struct {
	product     models.Product
	productList models.ProductList
	room    models.Room
}

func (handler HandleProductRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, user, _ := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user.Room = models.RoomByItsOwner(user.Id)

	json.NewDecoder(r.Body).Decode(&handler.product)

	handler.room = models.RoomByItsId(handler.product.ListRoom)
	handler.room.FindGuests()

	if !handler.room.IsOwner(user) {
		if !handler.room.IsGuest(user) {
			http.Error(w, "You are not a guest", http.StatusBadRequest)
			return
		}

		if handler.room.GuestPermission(user) < 2 {
			http.Error(w, "You do not have permission for this operation", http.StatusUnauthorized)
			return
		}
	}

	handler.productList.Name = handler.product.ListName
	handler.productList.Room = handler.product.ListRoom

	if err := models.InsertProduct(handler.product, handler.productList); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "Product name already in use", http.StatusConflict)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

  if _, ok := models.RoomProducts[handler.room.Id]; !ok {
    models.RoomProducts[handler.room.Id] = make(models.Products)
  }

  models.RoomProducts[handler.room.Id][handler.product.Name] = handler.product

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(handler.product)
}
