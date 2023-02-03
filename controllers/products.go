package controllers

import (
  "fmt"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/api/models"
)

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
  var both models.UpdatingProduct

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
