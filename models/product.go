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
	Description []byte  `json:"description"`
	Image       []byte  `json:"image"`
}

type UpdatingProduct struct {
	New Product `json:"new_product"`
	Old Product `json:"old_product"`
}

func InsertProduct(product Product) error {
	db := database.GetConnection()

  if product.ListName == "" {
    product.ListName = "orphans"
  }

	_, err := db.Query(database.InsertProduct,
		product.ListName, product.ListRoom, product.Name,
		product.Price, product.Description, product.Image)

	return err
}

func UpdateProduct(both UpdatingProduct, productList ProductList) error {
	db := database.GetConnection()

	_, err := db.Query(database.UpdateProduct,
		both.New.Name, both.New.Price, both.New.Description,
		both.New.Image, productList.Room, both.Old.Name)

	return err
}

func DeleteProduct(product Product) error {
	db := database.GetConnection()

	_, err := db.Query(database.DeleteProduct, product.Name, product.ListRoom)

	return err
}

func SelectOneProduct(room int, name string) (Product, error) {
	var p Product
	var db = database.GetConnection()

	rows, err := db.Query(database.SelectProduct, name, room)
	if err != nil {
		return p, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&p.ListName, &p.ListRoom, &p.Name, &p.Price, &p.Description, &p.Image)
		if err != nil {
			return Product{}, err
		}
	} else {
		return Product{}, errors.New(database.ErrNotfound)
	}

	return p, nil
}

func SelectProductByHisList(productName string, listRoom int) (error, Product) {
	var product Product
	var db = database.GetConnection()

	rows, err := db.Query(database.SelectProduct, productName, listRoom)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&product.ListName, &product.ListRoom, &product.Name, &product.Price, &product.Description, &product.Image)
		if err != nil {
			panic(err)
		}

		return nil, product
	}

	return errors.New("No product found"), product
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
