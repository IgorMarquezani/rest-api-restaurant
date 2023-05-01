package users

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/models"
)

func FullInfo(w http.ResponseWriter, r *http.Request) {
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

	for i := range user.RoomsAsGuest {
		user.RoomsAsGuest[i].FindGuests()
		user.RoomsAsGuest[i].FindProductsLists()

		for j := range user.RoomsAsGuest[i].ProductsList {
			user.RoomsAsGuest[i].ProductsList[j] = models.FindProductsInList(user.RoomsAsGuest[i].ProductsList[j])
		}

		user.RoomsAsGuest[i].FindTabs()
		user.RoomsAsGuest[i].FindTabsRequests()
	}

  user.ActiveRoom, _ = models.GetUserActiveRoom(user)
	user.ClearCriticalInfo()

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.Encode(user)
}
