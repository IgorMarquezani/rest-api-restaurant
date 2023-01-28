package models

type Product struct {
  ProductList ProductList `json:"product_list"`
  Name string `json:"name"`
  Price float64 `json:"price"`
  Description string `json:"description"`
  Image []byte `json:"image"`
}

type OldAndNew struct {
  New Product `json:"new_product"`
  Old Product `json:"old_product"`
}

