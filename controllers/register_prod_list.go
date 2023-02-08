package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/api/models"
)

type RegisterList struct {
}

func (rl RegisterList) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	var productList models.ProductList
  var jroom models.Room
  var err, u = verifySessionCookie(r)

  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  u.Room = models.SelectRoom(u.Id)
  u.Room.FindGuests()

	json.NewDecoder(r.Body).Decode(&productList)
  jroom = models.RoomById(productList.Room)
  jroom.FindGuests()

  if productList.Room != u.Room.Id {
    if jroom.Guests[u.Email].Id == 0 {
      http.Error(w, "You are not a guest in that room", http.StatusBadRequest)
      return
    }

    if models.SelectGuestPermission(u.Id, productList.Room) < 2 {
      http.Error(w, "You do not have permission", http.StatusBadRequest)
      return
    }
  }

	if err := models.InsertProductList(productList); err != nil {
		http.Error(w, "name already in use", http.StatusAlreadyReported)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
