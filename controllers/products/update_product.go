package products

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
	"github.com/gorilla/mux"
)

func Update(w http.ResponseWriter, r *http.Request) {
	var (
		oldProduct models.Product
		newProduct models.Product
		room       models.Room
	)

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := controllers.ValidJSONFormat(r.Body)
	if err != nil {
		http.Error(w, controllers.ErrJsonFormat, http.StatusBadRequest)
		return
	}

	json.Unmarshal(body, &newProduct)
	if newProduct.Name == "" {
		http.Error(w, models.ErrEmptyProductName, http.StatusExpectationFailed)
		return
	}

	if newProduct.ListRoom <= 0 {
		room = models.RoomByItsId(session.ActiveRoom)
		newProduct.ListRoom = room.Id
	} else {
		room = models.RoomByItsId(newProduct.ListRoom)
	}

	if !room.IsOwner(user) {
		if !room.IsGuest(user) {
			http.Error(w, models.ErrNotAGuest, http.StatusUnauthorized)
			return
		}

		if room.GuestPermission(user) < 2 {
			http.Error(w, models.ErrInsufficientPermission, http.StatusUnauthorized)
			return
		}
	}

	oldProduct, err = models.SelectOneProduct(room.Id, mux.Vars(r)["old-name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err, _ := models.SelectProductList(newProduct.ListName, room.Id); err != nil && newProduct.ListName != "" {
		models.InsertProductList(models.ProductList{
			Room: room.Id,
			Name: newProduct.ListName,
		})
	} 

	if err := models.UpdateProduct(newProduct, oldProduct, room.Id); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, models.ErrProductNameAlreadyUsed, http.StatusAlreadyReported)
			return
		}

		if database.IsForeignKeyConstraintError(err.Error()) {
			deleteErr := models.ProductErr{
				Title:    "Cannot update used(s) product(s)",
				Detail:   models.ErrProductStillUsed,
				Products: []string{oldProduct.Name},
			}

			w.WriteHeader(http.StatusConflict)
			controllers.EncodeJSON(w, deleteErr)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	controllers.EncodeJSON(w, newProduct)
}
