package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/api/models"
)

func NewProductList (w http.ResponseWriter, r *http.Request) {
  var productList models.ProductList
  json.NewDecoder(r.Body).Decode(&productList)
  
  if err := models.InsertProductList(productList); err != nil {
    http.Error(w, "Room does not exist or name already in use", http.StatusConflict)
    fmt.Println(err)
    return
  }

  w.WriteHeader(http.StatusCreated)
}

func NewInvite(w http.ResponseWriter, r *http.Request) {
  var invite models.Invite
  json.NewDecoder(r.Body).Decode(&invite)

  var owner = models.InitUserByRoom(invite.InvitingRoom) 
  owner.Room.FindGuests()

  var inviter = models.InitUserByCookie(r)

  if inviter.Email != owner.Email {
    if owner.Room.Guests[inviter.Email].Email == "" {
      http.Error(w, "You are not a guest in this room", http.StatusExpectationFailed)
      return
    }
  }

}

func NewGuest(w http.ResponseWriter, r *http.Request) {
  var guest models.Guest

  json.NewDecoder(r.Body).Decode(&guest)
}
