package models

import "github.com/api/database"

type PayedRequest struct {
	RoomId      uint   `json:"room_id"`
	TabId       uint64 `json:"tab_id"`
	ProductName string `json:"product_name"`
	Quantity    uint   `json:"quantity"`
}

func SelectPayedRequests(roomId uint, tabId uint64) ([]PayedRequest, error) {
	var (
		db            = database.GetConnection()
		payedRequests = make([]PayedRequest, 0)
	)

	sel, err := db.Query(database.SelectPayedRequests, roomId, tabId)
	if err != nil {
		return nil, err
	}

	defer sel.Close()

	for sel.Next() {
		pr := PayedRequest{}

		err := sel.Scan(&pr.RoomId, &pr.TabId, &pr.ProductName, &pr.Quantity)
		if err != nil {
			return nil, err
		}

		payedRequests = append(payedRequests, pr)
	}

	return payedRequests, nil
}

func NewPayedRequest(request Request) PayedRequest {
	return PayedRequest{
		RoomId:      uint(request.TabRoom),
		ProductName: request.ProductName,
		Quantity:    request.Quantity,
	}
}

func (pr PayedRequest) Insert() error {
	db := database.GetConnection()

	insert, err := db.Query(database.InsertPayedRequest,
		pr.RoomId, pr.TabId, pr.ProductName, pr.Quantity)

	if err != nil {
		return err
	}

	insert.Close()

	return nil
}
