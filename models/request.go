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
	Quantity        uint    `json:"quantity"`
}

type UpdatingRequest struct {
	TabRoom         int    `json:"tab_room"`
	TabNumber       int    `json:"tab_number"`
	ProductName     string `json:"product_name"`
	ProductListRoom int    `json:"product_list"`
	Quantity        uint   `json:"quantity"`

	Operation string `json:"operation"`
}

func InsertRequest(tab Tab, request Request) error {
	if request.ProductListRoom == 0 {
		request.ProductListRoom = tab.RoomId
	}

	if tab.RoomId != request.ProductListRoom {
		return errors.New("tab room and request room incompatibles")
	}

	var db = database.GetConnection()

	insert, err := db.Query(database.InsertRequest, tab.RoomId, tab.Number, request.ProductName, request.ProductListRoom, request.Quantity)
	if err != nil {
		return err
	}
	insert.Close()

	return nil
}

func SelectRequest(productName string, tabNumber, roomId int) (Request, error) {
	var request Request
	var db = database.GetConnection()

	rows, err := db.Query(database.SelectRequest, productName, tabNumber, roomId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&request.TabRoom, &request.TabNumber, &request.ProductName, &request.ProductListRoom, &request.Quantity)
    return request, err
	}

	return request, errors.New("No such request")
}

func UpdateRequestQuantity(request Request, quantity uint) error {
	var db = database.GetConnection()

	update, err := db.Query(database.UpdateRequestQuantity, quantity, request.ProductName, request.TabNumber, request.TabRoom)
	if err != nil {
		panic(err)
	}
	update.Close()

	return nil
}

func DeleteRequest(number, roomId int, productName string) {
	var db = database.GetConnection()

	del, err := db.Query(database.DeleteRequest, roomId, number, productName)
	if err != nil {
		panic(err)
	}
	del.Close()
}

func DeleteRequestsInTab(tab Tab) error {
	var db = database.GetConnection()

	del, err := db.Query(database.DeleteRequestsInTab, tab.RoomId, tab.Number)
	if err != nil {
		panic(err)
	}
	del.Close()

	return nil
}
