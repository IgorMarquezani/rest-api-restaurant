package models

import (
	_ "fmt"

	"github.com/api/database"
)

func InsertProductList (pr ProductList) error {
  db := database.GetConnection()

  _, err := db.Query(database.InsertProductList, pr.Name, pr.Room)

  return err
}

