package users

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/models"
)

type UserFullInfo struct {
}

func (ur UserFullInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, user := controllers.VerifySessionCookie(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusFailedDependency)
		return
	}

	user.Room = models.RoomByItsOwner(user.Id)
	user.Room.FindGuests()
	user.Room.FindProductsLists()

	// Finding all products in product list
	for i := range user.Room.ProductsList {
		user.Room.ProductsList[i] = models.FindProductsInList(user.Room.ProductsList[i])
	}

	user.Room.FindTabs()
	user.Room.FindTabsRequests()
	user.UserInvites()
	user.AcceptedRooms()
	user.ClearCriticalInfo()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
