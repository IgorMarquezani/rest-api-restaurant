package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/api/models"
)

type RegisterList struct {

}

func (rl RegisterList) verifyActiveRoom() {

}

func (rl RegisterList) verifyHashedPasswd() {

}

func ServeHTTP (w http.ResponseWriter, r *http.Request) {
  var productList models.ProductList
  json.NewDecoder(r.Body).Decode(&productList)
  
  if err := models.InsertProductList(productList); err != nil {
    http.Error(w, "Room does not exist or name already in use", http.StatusConflict)
    fmt.Println(err)
    return
  }

  w.WriteHeader(http.StatusCreated)
}
