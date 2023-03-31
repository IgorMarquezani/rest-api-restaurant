package models

import (
	"errors"

	"github.com/api/database"
)

type Products map[string]Product

type Product struct {
	ListName    string  `json:"list_name"`
	ListRoom    int     `json:"list_room"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       []byte  `json:"image"`
}

func (p *Product) Exists() bool {
	if p.ListRoom == 0 || p.ListName == "" || p.Name == "" {
		return false
	}

	err, _ := SelectProductByHisList(p.Name, p.ListRoom)
	if err != nil {
		return false
	}

	return true
}

type UpdatingProduct struct {
	New Product `json:"new_product"`
	Old Product `json:"old_product"`
}

func (up *UpdatingProduct) IsIncompatible() error {
	if up.New.ListRoom == 0 || up.Old.ListRoom == 0 {
		return errors.New("One of the products is missing his room id")
	}
	if up.New.ListName == "" || up.Old.ListName == "" {
		return errors.New("One of the products is missing his room name")
	}
	if up.New.ListName != up.Old.ListName {
		return errors.New("Incompatible lists names between products")
	}
	if up.New.ListRoom != up.Old.ListRoom {
		return errors.New("Incompatible rooms id between products")
	}

	return nil
}

func InsertProduct(product Product, productList ProductList) error {
	db := database.GetConnection()

	_, err := db.Query(
		database.InsertProduct,
		productList.Name, productList.Room,
		product.Name, product.Price, product.Description, product.Image)

	return err
}

func UpdateProduct(both UpdatingProduct, productList ProductList) error {
	db := database.GetConnection()

	_, err := db.Query(
		database.UpdateProduct,
		both.New.Name, both.New.Price, both.New.Description, both.New.Image,
		productList.Room, both.Old.Name)

	return err
}

func DeleteProduct(product Product) error {
	db := database.GetConnection()

	_, err := db.Query(
		database.DeleteProduct,
		product.Name, product.ListRoom)

	return err
}

func SelectProductByHisList(productName string, listRoom int) (error, Product) {
	var product Product
	var db = database.GetConnection()

	search, err := db.Query(
		database.SelectProduct, productName, listRoom)

	if err != nil {
		panic(err)
	}

	if search.Next() {
		err := search.Scan(
			&product.ListName, &product.ListRoom, &product.Name,
			&product.Price, &product.Description, &product.Image)

		if err != nil {
			panic(err)
		}

		return nil, product
	}

	return errors.New("No product found"), product
}
