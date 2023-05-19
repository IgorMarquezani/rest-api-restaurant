package models

import (
	"github.com/api/database"
)

type PayedTab struct {
	Id     uint64  `json:"id"`
	Number uint    `json:"number"`
	RoomId uint    `json:"room_id"`
	Value  float64 `json:"value"`
	Date   string  `json:"date"`
	Table  uint    `json:"table"`

	Requests []PayedRequest `json:"requests"`
}

func NewPayedTab(tab Tab, date string) PayedTab {
	pt := PayedTab{
		Number:   uint(tab.Number),
		RoomId:   uint(tab.RoomId),
		Value:    tab.PayValue,
		Date:     date,
		Table:    uint(tab.Table),
		Requests: make([]PayedRequest, len(tab.Requests)),
	}

	return pt
}

// this function should be called after inserting the payed_tab so it will have the Id field with some value
func (pt *PayedTab) AddRequests(requests []Request) {
	pt.Requests = make([]PayedRequest, len(requests))

	for _, request := range requests {
		pr := NewPayedRequest(request)

		pr.TabId = pt.Id
		pt.Requests = append(pt.Requests, pr)
	}
}

func (pt PayedTab) Insert() (uint64, error) {
	db := database.GetConnection()

	insert, err := db.Query(database.InsertPayedTab,
		pt.Number, pt.RoomId, pt.Value, pt.Date, pt.Table)

	if err != nil {
		return pt.Id, err
	}

	defer insert.Close()

	if insert.Next() {
		err = insert.Scan(&pt.Id)
	}

	return pt.Id, err
}

func (pt PayedTab) InsertRequests() error {
	for _, request := range pt.Requests {
		err := request.Insert()
		if err != nil {
			return err
		}
	}

	return nil
}

// func (pt PayedTab) InsertWithRequests() error {
// 	var err error

// 	pt.Id, err = pt.Insert()
// 	if err != nil {
// 		return err
// 	}

// 	err = pt.InsertRequests()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
