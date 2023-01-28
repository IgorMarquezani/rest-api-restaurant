package models

import (
	_ "fmt"

	"github.com/api/database"
)

func InsertUser(u User) error {
  db := database.GetConnection() 
  
  _, err := db.Query(database.InsertUser, u.Name, u.Email, u.Passwd, u.Img)

  return err 
}


func InsertProductList (pr ProductList) error {
  db := database.GetConnection()

  _, err := db.Query(database.InsertProductList, pr.Name, pr.Room)

  return err
}

func InsertProduct(product Product, productList ProductList) error {
  db := database.GetConnection()
  
  _, err := db.Query (
    database.InsertProduct,
    productList.Name, productList.Room,
    product.Name, product.Price, product.Description, product.Image)

  return err
}

func UpdateProduct (both OldAndNew, productList ProductList) error {
  db := database.GetConnection()

  _, err := db.Query(
    database.UpdateProduct, 
    both.New.Name, both.New.Price, both.New.Description, both.New.Image,
    productList.Room, both.Old.Name)

  return err 
}
