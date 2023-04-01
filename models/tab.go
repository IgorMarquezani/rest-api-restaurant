package models

import (
	"time"

	"github.com/api/database"
)

type Tab struct {
	Number   int     `json:"number"`
	RoomId   int     `json:"room"`
	PayValue float64 `json:"pay_value"`
	Maded    string  `json:"time_maded"`
	Table    int     `json:"table"`

	Requests []Request `json:"requests"`
}

func (t *Tab) CalculateValue() {
	for i := range t.Requests {
		product := RoomProducts[t.RoomId][t.Requests[i].ProductName]
		t.PayValue += product.Price * float64(t.Requests[i].Quantity)
	}
}

func (t *Tab) FindRequests() {
	var db = database.GetConnection()
	t.Requests = make([]Request, 0)

	search, err := db.Query(database.SelectRequestsInTab, t.RoomId, t.Number)
	if err != nil {
		panic(err)
	}

	for i := 0; search.Next(); i++ {
		t.Requests = append(t.Requests, Request{})

		search.Scan(&t.Requests[i].TabRoom, &t.Requests[i].TabNumber,
			&t.Requests[i].ProductName, &t.Requests[i].ProductListRoom, &t.Requests[i].Quantity)
	}
}

func NextTabNumberInRoom(room int) int {
	var (
		db       = database.GetConnection()
		selected Tab
		previous Tab
		next     Tab
	)

	search, err := db.Query(database.SelectTabsInRoom, room)
	if err != nil {
		panic(err)
	}

	for search.Next() {
		search.Scan(&next.Number, &next.RoomId, &next.PayValue, &next.Maded, &next.Table)

		if next.Number-previous.Number > 1 {
			selected.Number = previous.Number + 1
			break
		}

		previous = next
	}

	if selected.Number == 0 {
		return next.Number + 1
	}

	return selected.Number
}

func InsertTab(tab *Tab) error {
	var db = database.GetConnection()

	if tab.Number == 0 {
		tab.Number = NextTabNumberInRoom(tab.RoomId)
	}

	_, err := db.Query(database.InsertTab, tab.Number, tab.RoomId, time.Now().Local().Format("15:05:04"), tab.Table)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTab(tab Tab) error {
	var db = database.GetConnection()

	_, err := db.Query(database.DeleteTab, tab.Number, tab.RoomId)
	if err != nil {
		panic(err)
	}

	return nil
}
