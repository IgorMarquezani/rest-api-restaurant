package prodlist

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/models"
)

type RegisterList struct {
  productList models.ProductList
  jsonRoom    models.Room
}

func (rl RegisterList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err, user = controllers.VerifySessionCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Room = models.RoomByItsOwner(user.Id)

	json.NewDecoder(r.Body).Decode(&rl.productList)
	rl.jsonRoom = models.RoomByItsId(rl.productList.Room)
	rl.jsonRoom.FindGuests()

	if !rl.jsonRoom.IsOwner(user) {
		if !rl.jsonRoom.IsGuest(user) {
			http.Error(w, "You are not a guest in that room", http.StatusBadRequest)
			return
		}

		if rl.jsonRoom.GuestPermission(user) < 2 {
			http.Error(w, "You do not have permission", http.StatusBadRequest)
			return
		}
	}

	if err := models.InsertProductList(rl.productList); err != nil {
		http.Error(w, "name already in use", http.StatusAlreadyReported)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
