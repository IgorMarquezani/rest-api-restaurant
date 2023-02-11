package products

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/models"
)

type DeleteProduct struct {
	product  models.Product
	jsonRoom models.Room
}

func (dp DeleteProduct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, user, _ := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user.Room = models.RoomByItsOwner(user.Id)

	json.NewDecoder(r.Body).Decode(&dp.product)

	dp.jsonRoom = models.RoomByItsId(dp.product.ListRoom)
	dp.jsonRoom.FindGuests()

	if !dp.jsonRoom.IsOwner(user) {
		if !dp.jsonRoom.IsGuest(user) {
			http.Error(w, "Not a guest in that room", http.StatusUnauthorized)
			return
		}

		if dp.jsonRoom.GuestPermission(user) < 2 {
			http.Error(w, "Don't have permission for this operation", http.StatusUnauthorized)
			return
		}
	}

	if err := models.DeleteProduct(dp.product); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dp.product)
}
