package users

import (
	"net/http"

	"github.com/api/controllers"
	"github.com/api/models"
)

type UserRooms struct {
  
}

func (ur UserRooms) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  err, user := controllers.VerifySessionCookie(r)
  if err != nil {
    http.Error(w, "Not logged in", http.StatusFailedDependency)
    return
  }

  user.Room = models.RoomByItsOwner(user.Id)
  user.Room.FindTabs()


}
