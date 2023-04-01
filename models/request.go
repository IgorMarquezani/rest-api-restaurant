package models

import (
	"errors"

	"github.com/api/database"
)

type Request struct {
	TabRoom         int    `json:"tab_room"`
	TabNumber       int    `json:"tab_number"`
	ProductName     string `json:"product_name"`
	ProductListRoom int    `json:"product_list"`
	Quantity        int    `json:"quantity"`
}

func InsertRequest(tab Tab, request Request) error {
	if request.ProductListRoom == 0 {
		request.ProductListRoom = tab.RoomId
	}

	if tab.RoomId != request.ProductListRoom {
		return errors.New("tab room and request room incompatibles")
	}

	var db = database.GetConnection()

	_, err := db.Query(database.InsertRequest, tab.RoomId, tab.Number, request.ProductName, request.ProductListRoom, request.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRequestsInTab(tab Tab) error {
	var db = database.GetConnection()

	_, err := db.Query(database.DeleteRequestsInTab, tab.RoomId, tab.Number)
	if err != nil {
		panic(err)
	}

	return nil
}
