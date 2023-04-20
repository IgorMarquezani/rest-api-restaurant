package models

import (
	"errors"
	"time"

	"github.com/api/database"
)

const (
  // No such tab in this room
	ErrNoSuchTab        string = "No such tab in this room"
  // Invalid number for request
	ErrInvalidTabNumber string = "Invalid number for request"
)

type Tab struct {
	Number   int     `json:"number"`
	RoomId   int     `json:"room"`
	PayValue float64 `json:"pay_value"`
	Maded    string  `json:"time_maded"`
	Table    int     `json:"table"`

	Requests []Request `json:"requests"`
}

type UpdatingTab struct {
	Number   int     `json:"number"`
	RoomId   int     `json:"room"`
	PayValue float64 `json:"pay_value"`
	Maded    string  `json:"time_maded"`
	Table    int     `json:"table"`

	Requests []UpdatingRequest `json:"requests"`
}

func SelectTabByNumber(number, roomId int) (Tab, error) {
	var db = database.GetConnection()
	var tab Tab

	rows, err := db.Query(database.SelectTab, number, roomId)
	if err != nil {
		panic(err)
	}
  defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&tab.Number, &tab.RoomId, &tab.PayValue, &tab.Maded, &tab.Table)
		if err != nil {
			panic(err)
		}

		return tab, nil
	}

	tab.FindRequests()

	return tab, errors.New(ErrNoSuchTab)
}

func InsertTab(tab *Tab) error {
	var db = database.GetConnection()

	if tab.Number == 0 {
		tab.Number = NextTabNumberInRoom(tab.RoomId)
	}

	if tab.PayValue == 0 {
		tab.CalculateValue()
	}

  if _, err := time.Parse("15:04:05", tab.Maded); err != nil {
	  tab.Maded = time.Now().Local().Format("15:04:05")
  }

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

func UpdateTab(oldTab, newTab Tab) error {
	var db = database.GetConnection()

	_, err := db.Query(database.UpdateTab, newTab.Table, oldTab.Number, oldTab.RoomId)
	if err != nil {
		return err
	}

	return nil
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

func (t *Tab) SortRequest() []Request {
	if t.Requests == nil {
		t.FindRequests()
	}

	requests := t.Requests

	for i := range requests {
		for j := i; j < len(requests); j++ {
			if requests[j].ProductName > requests[i].ProductName {
				tmp := requests[j]
				requests[j] = requests[i]
				requests[i] = tmp
			}
		}
	}

	t.Requests = requests

	return requests
}

func NextTabNumberInRoom(room int) int {
	var (
		selected Tab
		previous Tab
		next     Tab
		db       = database.GetConnection()
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
