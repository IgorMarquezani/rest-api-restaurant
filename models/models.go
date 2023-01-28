package models

import (
	"database/sql"
	_ "fmt"

	"github.com/api/database"
)

func InsertUser(u User) error {
  db := database.GetConnection() 
  
  _, err := db.Query(database.InsertUser, u.Name, u.Email, u.Passwd, u.Img)

  return err 
}

func SelectUser(key any) User {
  var search *sql.Rows
  var err error

  db := database.GetConnection()
  u := User{}

  switch data := key.(type) {
  case string:
    search, err = db.Query(database.SearchUserByEmail, data)
  case int:
    search, err = db.Query(database.SearchUserById, data)
  default:
    panic("Invalid data type")
  }

  if err != nil {
    panic(err)
  }

  if search.Next() {
    err = search.Scan(&u.Id, &u.Name, &u.Email, &u.Passwd, &u.Img)
    if err != nil {
      panic(err)
    }
  }

  return u
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
