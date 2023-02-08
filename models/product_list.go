package models

import "github.com/api/database"

type ProductListName string

type ProdListMap map[ProductListName]ProductList

type ProductList struct {
	Name string `json:"name"`
	Room int    `json:"room"`
  Products ProductMap
}

func InsertProductList(pr ProductList) error {
	db := database.GetConnection()

	_, err := db.Query(database.InsertProductList, pr.Name, pr.Room)

	return err
}
