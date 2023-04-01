package main

import (
	"fmt"

	"github.com/api/database"
	"github.com/api/models"
	"github.com/api/routes"
)

func init() {
	database.MustNewConnection()

	var (
		db      = database.GetConnection()
		room    models.Room
		product models.Product
	)

	models.RoomProducts = make(models.ProductsInRoom)

	searchR, err := db.Query(database.SelectAllRooms)
	if err != nil {
		panic(err)
	}
	defer searchR.Close()

	for searchR.Next() {
		searchR.Scan(&room.Id, &room.Owner)
		searchP, err := db.Query(database.SelectProductsByRoom, room.Id)
		if err != nil {
			panic(err)
		}

		for searchP.Next() {
			var (
				description []byte
				img         []byte
			)

			err := searchP.Scan(&product.ListName, &product.ListRoom, &product.Name, &product.Price, &description, &img)
			if err != nil {
				panic(err)
			}

			if description != nil {
				product.Description = string(description)
			}

			if img != nil {
				product.Image = img
			}

			models.RoomProducts[room.Id] = make(models.Products)
			models.RoomProducts[room.Id][product.Name] = product

		}

		searchP.Close()
	}

	for i := range models.RoomProducts {
		for j := range models.RoomProducts[i] {
			fmt.Println(models.RoomProducts[i][j])
		}
	}
}

func main() {
	fmt.Println("Initializing API for Fatec project")
	routes.HandleRequest()
}
