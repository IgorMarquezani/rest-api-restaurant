package models

import (
	"github.com/api/database"
)

type Tab struct {
	Number   int     `json:"number"`
	RoomId   int     `json:"room"`
	PayValue float64 `json:"pay_value"`
	Maded    string  `json:"time_maded"`

	Requests []Request `json:"requests"`
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
			&t.Requests[i].ProductName, &t.Requests[i].ProductList, &t.Requests[i].Quantity)
	}
}

func NextTabNumberInRoom(room int) int {
	var db = database.GetConnection()
	var selected Tab
	var previous Tab
	var next Tab

	search, err := db.Query(database.SelectTabsInRoom, room)
	if err != nil {
		panic(err)
	}

	for search.Next() {
		search.Scan(&next.Number, &next.RoomId, &next.PayValue, &next.Maded)

		if next.Number-previous.Number > 1 {
			selected.Number = previous.Number + 1
			break
		}

		previous = next
	}

	if selected.Number == 0 {
		max, _ := db.Query(database.SelectMaxTabId, room)
		if max.Next() {
			max.Scan(&selected.Number)
		}
	}

	return selected.Number
}

func InsertTab(tab *Tab) error {
	var db = database.GetConnection()

	if tab.Number == 0 {
		tab.Number = NextTabNumberInRoom(tab.RoomId)
	}

	insert, err := db.Query(database.InsertTab, tab.Number, tab.RoomId)
	if err != nil {
		return err
	}

	if insert.Next() {
		insert.Scan(tab.Number, tab.RoomId, &tab.PayValue, &tab.Maded)
	}

	return nil
}
