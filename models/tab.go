package models

import (
	"errors"
	"sort"
	"time"

	"github.com/api/database"
)

const (
	// No such tab in this room
	ErrNoSuchTab string = "No such tab in this room"
	// Invalid number for request
	ErrInvalidTabNumber string = "Invalid number for request"
)

type Tab struct {
	RoomId   int     `json:"room"`
	Number   int     `json:"number"`
	PayValue float64 `json:"pay_value"`
	Maded    string  `json:"time_maded"`
	Table    int     `json:"table"`

	Requests []Request `json:"requests"`
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

		tab.FindRequests()

		return tab, nil
	}

	return tab, errors.New(ErrNoSuchTab)
}

func InsertTab(tab *Tab) error {
	var db = database.GetConnection()

	if tab.PayValue == 0 {
		tab.CalculateValue()
	}

	if _, err := time.Parse("15:04:05", tab.Maded); err != nil {
		tab.Maded = time.Now().Local().Format("15:04:05")
	}

	rows, err := db.Query(database.InsertTab, tab.RoomId, tab.RoomId, tab.PayValue, tab.Maded, tab.Table)
	if err != nil {
		return err
	}

  if rows.Next() {
    if err := rows.Scan(&tab.Number); err != nil {
      return err
    }
  }

	rows.Close()

	return nil
}

func DeleteTab(tab Tab) error {
	var db = database.GetConnection()

	rows, err := db.Query(database.DeleteTab, tab.Number, tab.RoomId)
	if err != nil {
		panic(err)
	}
	rows.Close()

	return nil
}

func UpdateTab(oldTab, newTab Tab) error {
	var db = database.GetConnection()

	rows, err := db.Query(database.UpdateTab, newTab.Table, oldTab.Number, oldTab.RoomId)
	if err != nil {
		return err
	}
	rows.Close()

	return nil
}

func IncreaseTabValue(tabNumber, roomId int, value float64) error {
	db := database.GetConnection()

	update, err := db.Query(database.IncreaseTabValue, tabNumber, roomId, value, tabNumber, roomId)
	if err != nil {
		return err
	}
	update.Close()

	return nil
}

func DecreaseTabValue(tabNumber, roomId int, value float64) error {
	db := database.GetConnection()

	del, err := db.Query(database.DecreaseTabValue, tabNumber, roomId, value, tabNumber, roomId)
	if err != nil {
		return err
	}
	del.Close()

	return nil
}

func (t *Tab) Len() int {
	return len(t.Requests)
}

func (t *Tab) Less(i, j int) bool {
	return t.Requests[i].ProductName < t.Requests[j].ProductName
}

func (t *Tab) Swap(i, j int) {
	request := t.Requests[i]
	t.Requests[i] = t.Requests[j]
	t.Requests[j] = request
}

func (t *Tab) GroupRequests() {
	var (
		total   uint
		grouped []Request
	)

	sort.Sort(t)

	for i := 0; i < t.Len(); i++ {
		request := t.Requests[i]

		for ; i < t.Len()-1 && t.Requests[i].ProductName == t.Requests[i+1].ProductName; i++ {
			request.Quantity += t.Requests[i+1].Quantity
		}

		grouped = append(grouped, request)
		total++
	}

	t.Requests = grouped[0:total]
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

// func (t *Tab) SortRequest() []Request {
// 	if t.Requests == nil {
// 		t.FindRequests()
// 	}

// 	requests := t.Requests

// 	for i := range requests {
// 		for j := i; j < len(requests); j++ {
// 			if requests[j].ProductName > requests[i].ProductName {
// 				tmp := requests[j]
// 				requests[j] = requests[i]
// 				requests[i] = tmp
// 			}
// 		}
// 	}

// 	t.Requests = requests

// 	return requests
// }

func (tab *Tab) RemoveMadedTrash() error {
	if len(tab.Maded) < 20 {
		return errors.New("Too short")
	}

	tab.Maded = tab.Maded[11:19]

	return nil
}
