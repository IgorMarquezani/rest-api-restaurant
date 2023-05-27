package models

import "github.com/api/database"

type PayedTab struct {
	Id     uint64  `json:"id"`
	Number uint    `json:"number"`
	RoomId uint    `json:"room_id"`
	Value  float64 `json:"value"`
	Date   string  `json:"date"`
	Table  uint    `json:"table"`

	Requests []PayedRequest `json:"requests"`
}

func SelectPayedTabs(from, to string) ([]PayedTab, error) {
	var (
		db        = database.GetConnection()
		payedTabs = make([]PayedTab, 0)
	)

	sel, err := db.Query(database.SelectPayedTabWithInterval, from, to)
	if err != nil {
		return nil, err
	}

	defer sel.Close()

	for sel.Next() {
		pt := PayedTab{}

		err := sel.Scan(&pt.Id, &pt.Number, &pt.RoomId, &pt.Value, &pt.Date, &pt.Table)
		if err != nil {
			return nil, err
		}

		payedTabs = append(payedTabs, pt)
	}

	return payedTabs, nil
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

func NewPayedTab(tab Tab, date string) PayedTab {
	return PayedTab{
		Number:   uint(tab.Number),
		RoomId:   uint(tab.RoomId),
		Value:    tab.PayValue,
		Date:     date,
		Table:    uint(tab.Table),
		Requests: make([]PayedRequest, len(tab.Requests)),
	}
}

func (pt *PayedTab) AddRequests(requests []Request) {
	pt.Requests = make([]PayedRequest, len(requests))

	for i, request := range requests {
		pr := NewPayedRequest(request)

		pr.TabId = pt.Id
		pt.Requests[i] = pr
	}
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
