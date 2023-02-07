package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/api/models"
)

type RegisterProduct struct {
  product models.Product
}

func (rp RegisterProduct) NewProduct (w http.ResponseWriter, r *http.Request) {
	json.NewDecoder(r.Body).Decode(&rp.product)

	if err := models.InsertProduct(rp.product, rp.product.ProductList); err != nil {
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
