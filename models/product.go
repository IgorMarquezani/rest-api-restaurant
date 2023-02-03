package models

import(
  "github.com/api/database"
)

type Product struct {
  ProductList ProductList `json:"product_list"`
  Name string `json:"name"`
  Price float64 `json:"price"`
  Description string `json:"description"`
  Image []byte `json:"image"`
}

type UpdatingProduct struct {
  New Product `json:"new_product"`
  Old Product `json:"old_product"`
}

func InsertProduct(product Product, productList ProductList) error {
  db := database.GetConnection()
  
  _, err := db.Query (
    database.InsertProduct,
    productList.Name, productList.Room,
    product.Name, product.Price, product.Description, product.Image)

  return err
}

func UpdateProduct (both UpdatingProduct, productList ProductList) error {
  db := database.GetConnection()

  _, err := db.Query(
    database.UpdateProduct, 
    both.New.Name, both.New.Price, both.New.Description, both.New.Image,
    productList.Room, both.Old.Name)

  return err 
}
