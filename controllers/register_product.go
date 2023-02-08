package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/api/models"
)

type RegisterProduct struct {
	product     models.Product
	productList models.ProductList
	jsonRoom    models.Room
}

func (rp RegisterProduct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err, user = verifySessionCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	user.Room = models.SelectRoom(user.Id)
	user.Room.FindGuests()

	json.NewDecoder(r.Body).Decode(&rp.product)
	rp.productList.Name = rp.product.ListName
	rp.productList.Room = rp.product.ListRoom
	rp.jsonRoom = models.RoomById(rp.product.ListRoom)
	rp.jsonRoom.FindGuests()

	if user.Room.Id != rp.jsonRoom.Id {
		if rp.jsonRoom.Guests[user.Email].Id == 0 {
			http.Error(w, "You are not a guest", http.StatusBadRequest)
			return
		}

		http.Error(w, "You are not the owner of that room", http.StatusBadRequest)
		return
	}

	if err := models.InsertProduct(rp.product, rp.productList); err != nil {
		message := strings.Split(err.Error(), " ")
		if message[1] == "duplicate" && message[2] == "key" {
			http.Error(w, "Product name already in use", http.StatusConflict)
			return
		}

		http.Error(w, "Unknow error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
