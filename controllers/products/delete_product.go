package products

import (
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
	"github.com/gorilla/mux"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	var (
		product models.Product
		room    models.Room
	)

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	productName := mux.Vars(r)["name"]
	if productName == "" {
		http.Error(w, models.ErrEmptyProductName, http.StatusBadRequest)
		return
	}

	room = models.RoomByItsId(session.ActiveRoom)

	product, err = models.SelectOneProduct(room.Id, productName)
	if err != nil {
		http.Error(w, models.ErrNoSuchProduct, http.StatusNoContent)
		return
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

	if err := models.DeleteProduct(product); err != nil {
		if database.IsForeignKeyConstraintError(err.Error()) {
			deleteErr := models.ProductErr{
				Title:    "Cannot delete used(s) product(s)",
				Detail:   models.ErrProductStillUsed,
				Products: []string{product.Name},
			}

			w.WriteHeader(http.StatusConflict)
			controllers.EncodeJSON(w, deleteErr)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
