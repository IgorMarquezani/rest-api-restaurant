package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/api/database"
	"github.com/api/models"
	"github.com/api/routes"
)

func init() {
	database.MustNewConnection()

	var db = database.GetConnection()
	models.RoomProducts = make(models.ProductsInRoom)

	searchR, err := db.Query(database.SelectAllRooms)
	if err != nil {
		if database.IsAuthenticantionFailedError(err.Error()) {
			log.Println("error: ", err.Error(), "\nAre you shure the database is up?")
			os.Exit(1)
		}
		panic(err)
	}
	defer searchR.Close()

	for searchR.Next() {
		var (
			room models.Room
		)

		searchR.Scan(&room.Id, &room.Owner)

		models.RoomProducts[room.Id] = make(models.Products)

		searchP, err := db.Query(database.SelectProductsByRoom, room.Id)
		if err != nil {
			panic(err)
		}
		defer searchP.Close()

		for searchP.Next() {
			var (
				product models.Product
				desc    []byte
				img     []byte
			)
			err := searchP.Scan(&product.ListName, &product.ListRoom, &product.Name, &product.Price, &desc, &img)
			if err != nil {
				panic(err)
			}

			if desc != nil {
				product.Description = desc
			}

			if img != nil {
				product.Image = img
			}

			models.RoomProducts[room.Id][product.Name] = product
		}
	}

	for i := range models.RoomProducts {
		for j := range models.RoomProducts[i] {
			fmt.Println(models.RoomProducts[i][j])
		}
	}
}

func main() {
	var (
		max int = runtime.NumCPU()
		err error
	)

	if len(os.Args) > 1 {
		max, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Invalid number for max threads")
		} else {
			runtime.GOMAXPROCS(max)
		}
	}

	fmt.Println("Initializing API for Fatec project")
	fmt.Printf("Running with a total of %d threads available\n", max)
	routes.HandleRequest()
}
