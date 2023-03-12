package main

import (
	"fmt"

	"github.com/api/database"
	"github.com/api/routes"
)

func main() {
	database.NewConnection()
	fmt.Println("Initializing API for Fatec project")
	routes.HandleRequest()
}
