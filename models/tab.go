package models

import (
	"errors"
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

func InsertTab(tab *Tab) error {
	var db = database.GetConnection()

	if tab.Number == 0 {
		tab.Number = NextTabNumberInRoom(tab.RoomId)
	}

	if tab.PayValue == 0 {
		tab.CalculateValue()
	}

  tab.Maded = time.Now().Local().Format("15:04:05")

	_, err := db.Query(database.InsertTab, tab.Number, tab.RoomId, tab.PayValue, tab.Maded, tab.Table)
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

func (t *Tab) CalculateValue() {
	for i := range t.Requests {
		p := RoomProducts[t.RoomId][t.Requests[i].ProductName]
		t.PayValue += p.Price * float64(t.Requests[i].Quantity)
	}
}

func (t *Tab) FindRequests() {
	var db = database.GetConnection()
	t.Requests = make([]Request, 0)

	rows, err := db.Query(database.SelectRequestsInTab, t.RoomId, t.Number)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		t.Requests = append(t.Requests, Request{})

		rows.Scan(&t.Requests[i].TabRoom, &t.Requests[i].TabNumber,
			&t.Requests[i].ProductName, &t.Requests[i].ProductListRoom, &t.Requests[i].Quantity)
	}
}

func NextTabNumberInRoom(room int) int {
	var db = database.GetConnection()
	var (
		selected Tab
		previous Tab
		next     Tab
	)

	rows, err := db.Query(database.SelectTabsInRoom, room)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&next.Number, &next.RoomId, &next.PayValue, &next.Maded, &next.Table)

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

func (tab *Tab) RemoveMadedTrash() error {
	if len(tab.Maded) < 20 {
		return errors.New("Too short")
	}

	tab.Maded = tab.Maded[11:19]

	return nil
}
