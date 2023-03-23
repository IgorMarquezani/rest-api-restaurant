package rooms

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/api/controllers"
	"github.com/api/models"
)

type RoomInfo struct {

}

func (ri RoomInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  err, user := controllers.VerifySessionCookie(r)
  if err != nil {
    http.Error(w, "Not loged in", http.StatusFailedDependency)
    return
  }

  roomStr := r.URL.Query().Get("id")
  roomId, _ := strconv.Atoi(roomStr)
  room := models.RoomByItsId(roomId)

  if !room.IsOwner(user) {
    if !room.IsGuest(user) {
      http.Error(w, "Not owner or a guest", http.StatusUnauthorized)
      return
    }
  }

  room.FindGuests()
  room.FindProductsLists()

  for i := range room.ProductsList {
    room.ProductsList[i] = models.FindProductsInList(room.ProductsList[i])
  }

  room.FindTabs()
  room.FindTabsRequests()

  json.NewEncoder(w).Encode(room)
}
