package users

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/models"
)

func WithActiveRoom(w http.ResponseWriter, r *http.Request)  {
	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusFailedDependency)
		return
	}

  user.Room = models.RoomByItsId(session.ActiveRoom)

	user.Room.FindGuests()
	user.Room.FindProductsLists()

	// Finding all products in product list
	for i := range user.Room.ProductsList {
		user.Room.ProductsList[i] = models.FindProductsInList(user.Room.ProductsList[i])
	}

	user.Room.FindTabs()
	user.Room.FindTabsRequests()
	user.UserInvites()

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.Encode(user)
}
