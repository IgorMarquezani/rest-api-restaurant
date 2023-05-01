package models

import (
	"errors"

	"github.com/api/database"
)

const (
	ErrNoSuchProduct          = "No such product in this room"
	ErrProductNameAlreadyUsed = "Name already used in this room"
	ErrEmptyProductName       = "Product Name is empty"
	ErrProductStillUsed       = "There is one or more request on some tab(s) that is using this(e) product(s)"
)

type Products map[string]Product

type ProductErr struct {
	Title    string   `json:"title"`
	Detail   string   `json:"detail"`
	Products []string `json:"products"`
}

type Product struct {
	ListName    string  `json:"list_name"`
	ListRoom    int     `json:"list_room"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       []byte  `json:"image"`
}

func InsertProduct(product Product) error {
	db := database.GetConnection()

	if product.ListName == "" {
		product.ListName = "orphans"
	}

	insert, err := db.Query(database.InsertProduct,
		product.ListName, product.ListRoom, product.Name,
		product.Price, product.Description, product.Image)

	if err != nil {
		return err
	}

	insert.Close()

	return nil
}

func UpdateProduct(New, Old Product, roomId int) error {
	db := database.GetConnection()

	update, err := db.Query(database.UpdateProduct,
		New.Name, New.Price, New.Description,
		New.Image, roomId, Old.Name)

	if err != nil {
		return err
	}

	update.Close()

	return nil
}

func DeleteProduct(product Product) error {
	db := database.GetConnection()

	del, err := db.Query(database.DeleteProduct, product.Name, product.ListRoom)
	if err != nil {
		return err
	}

	del.Close()

	return nil
}

func SelectOneProduct(room int, name string) (Product, error) {
	var (
		product Product
		db      = database.GetConnection()
	)

	rows, err := db.Query(database.SelectProduct, name, room)
	if err != nil {
		return product, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&product.ListName, &product.ListRoom, &product.Name, &product.Price, &product.Description, &product.Image)
		return product, err
	}

	return product, errors.New(ErrNoSuchProduct)
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
