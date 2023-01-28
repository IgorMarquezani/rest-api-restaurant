package controllers

import (
	"encoding/json"
	"fmt"
	_ "io"
	"net/http"
	_ "strconv"
	"strings"

	"github.com/api/models"
)

func NewUser(w http.ResponseWriter, r *http.Request) {
  var u models.User
  json.NewDecoder(r.Body).Decode(&u)

  if err := models.InsertUser(u); err != nil {
    fmt.Println(err)
    http.Error(w, "E-mail already in use", http.StatusConflict)
    return
  }

  w.WriteHeader(http.StatusCreated)
}

func EnterUser(w http.ResponseWriter, r *http.Request) {
  var u models.User
  json.NewDecoder(r.Body).Decode(&u)
  ur := models.SelectUser(u.Email);

  if u.Email != ur.Email {
    http.Error(w, "E-mail not registered", http.StatusNotFound)
    return
  }

  if u.Passwd != ur.Passwd {
    http.Error(w, "Incompatible password", http.StatusConflict)
    return
  }

  w.WriteHeader(http.StatusAccepted)
  json.NewEncoder(w).Encode(ur)
}

func UpdateUser (w http.ResponseWriter, r *http.Request) {
  /*
  var u, ur models.User
  json.NewDecoder(r.Body).Decode(&u)
  ur = models.SelectUser(u.Id)
  */
}

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

func NewProduct(w http.ResponseWriter, r *http.Request) {
  var product models.Product

  json.NewDecoder(r.Body).Decode(&product)

  if err := models.InsertProduct(product, product.ProductList); err != nil {
    message := strings.Split(err.Error(), " ")

    if message[1] == "duplicate" && message[2] == "key" {
      http.Error(w, "Product Already Registered", http.StatusAlreadyReported)
      return
    }

    http.Error(w, "Not Know Error", http.StatusConflict)
    fmt.Println(err)
    return
  }

  w.WriteHeader(http.StatusCreated)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
  var both models.OldAndNew

  json.NewDecoder(r.Body).Decode(&both)

  if both.Old.ProductList.Name == "" || both.Old.ProductList.Room == 0 {
    http.Error(w, "Missing product_list information", http.StatusBadRequest)
    return
  }

  if err := models.UpdateProduct(both, both.Old.ProductList); err != nil {
    http.Error(w, "Unknow error", http.StatusConflict)
    message := strings.Split(err.Error(), " ")
    fmt.Println(message)
    return
  }

  w.WriteHeader(http.StatusAccepted)
}

func NewInvite(w http.ResponseWriter, r *http.Request) {
  var user = models.InitUserByCookie(r)
  var invite models.Invite

  json.NewDecoder(r.Body).Decode(&invite)

  if user.Room.Id != invite.InvitingRoom {
    http.Error(w, "Cannot invite to a room that you dont have permission", http.StatusConflict)
    return
  }
}

func NewGuest(w http.ResponseWriter, r *http.Request) {
  var guest models.Guest

  json.NewDecoder(r.Body).Decode(&guest)

}
