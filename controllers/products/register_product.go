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
	if product.Name == "" {
		http.Error(w, "Cannot register a product with empty name", http.StatusBadRequest)
		return
	}

	if product.ListRoom <= 0 {
		room = models.RoomByItsId(session.ActiveRoom)
		product.ListRoom = room.Id
	} else {
		room = models.RoomByItsId(product.ListRoom)
	}

  ok, permission := room.IsOwnerOrGuest(user)
  if !ok {
			http.Error(w, models.ErrNotAGuest, http.StatusBadRequest)
			return
  }

  if permission < 2 {
			http.Error(w, models.ErrInsufficientPermission, http.StatusUnauthorized)
			return
  }

	if product.ListName == "" {
		product.ListName = "orphans"
	}

	if err := models.InsertProduct(product); err != nil {
		if database.IsForeignKeyConstraintError(err.Error()) {
			list := models.ProductList{
				Room: product.ListRoom,
				Name: product.ListName,
			}

			models.InsertProductList(list)
			models.InsertProduct(product)

			w.WriteHeader(http.StatusCreated)
			controllers.EncodeJSON(w, product)
      return
		}

		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, models.ErrProductNameAlreadyUsed, http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	controllers.EncodeJSON(w, product)
}
