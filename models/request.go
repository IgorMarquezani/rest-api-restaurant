package models

type Request struct {
	TabRoom     int    `json:"tab_room"`
	TabNumber   int    `json:"tab_number"`
	ProductName string `json:"product_name"`
	ProductList string `json:"product_list"`
	Quantity    int    `json:"quantity"`
}
