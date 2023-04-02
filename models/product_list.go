package models

import (
	"errors"

	"github.com/api/database"
)

type ProductListMap map[string]ProductList

type ProductList struct {
	Name     string `json:"name"`
	Room     int    `json:"room"`
	Products Products
}

func InsertProductList(pr ProductList) error {
	db := database.GetConnection()

	_, err := db.Query(database.InsertProductList, pr.Name, pr.Room)

	return err
}

func SelectProductList(name string, roomId int) (error, ProductList) {
	var prodList ProductList
	var db = database.GetConnection()

	rows, err := db.Query(database.SelectProductList, name, roomId)
	if err != nil {
		panic(err)
	}
  defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&prodList.Name, &prodList.Room)
		if err != nil {
			panic(err)
		}

		return nil, prodList
	}

	return errors.New("No such product list"), prodList
}

func FindProductsInList(pl ProductList) ProductList {
	var db = database.GetConnection()

	pl.Products = make(Products)

	rows, err := db.Query(database.SelectProductsByList, pl.Name, pl.Room)
	if err != nil {
		panic(err)
	}
  defer rows.Close()

	for rows.Next() {
		product := Product{}

		rows.Scan(&product.ListName, &product.ListRoom,
			&product.Name, &product.Price, &product.Description, &product.Image)

		pl.Products[product.Name] = product
	}

	return pl
}

func (pl *ProductList) Exists() bool {
	if pl.Name == "" || pl.Room == 0 {
		return false
	}

	err, _ := SelectProductList(pl.Name, pl.Room)
	if err != nil {
		return false
	}

	return true
}

func (pl *ProductList) IsOnProductList(product Product) bool {
	if pl.Products == nil {
		pl.GetProducts()
	}

	if pl.Products[product.Name].Name != "" {
		return true
	}

	return false
}

func (pl *ProductList) GetProducts() Products {
	var (
	  db = database.GetConnection() 
    product Product
  )

	if pl.Products == nil {
		pl.Products = make(Products)
	}

	rows, err := db.Query(database.SelectProductsByList, pl.Name, pl.Room)
	if err != nil {
		panic(err)
	}
  defer rows.Close()

	for rows.Next() {
		rows.Scan(&product.ListName, &product.ListRoom,
			&product.Name, &product.Price, &product.Description, &product.Image)

		pl.Products[product.Name] = product
	}

	return pl.Products
}

