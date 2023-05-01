package models

import (
	"github.com/api/database"
)

type RoomsProductsChan map[int]chan Product

var roomsProductsChan = make(RoomsProductsChan)

type ProductsInRoom map[int]Products

var RoomProducts = make(ProductsInRoom)

type Room struct {
	Id    int `json:"id"`
	Owner int `json:"owner"`

	Guests       GuestMap
	ProductsList ProductListMap
	Tabs         []Tab
}

const (
	ErrInsufficientPermission = "Does not have permission for this operation"
)

func CacheProduct(roomId int, product Product) {

}

func AddProductToRoomCache(roomId int) {
	var c chan Product
	var ok bool

	if c, ok = roomsProductsChan[roomId]; !ok {
		roomsProductsChan[roomId] = make(chan Product, 10)
		c = roomsProductsChan[roomId]
	}

	go func() {
		if _, ok := RoomProducts[roomId]; !ok {
			RoomProducts[roomId] = make(Products)
		}

		for {
			select {
			case product := <-c:
				RoomProducts[roomId][product.Name] = product
			}
		}
	}()
}

// Working everything down here

func MustInsertRoom(userId uint) {
	db := database.GetConnection()

	insert, err := db.Query(database.InsertRoom, userId)
	if err != nil {
		panic(err)
	}
	insert.Close()
}

/*
Don't ever call this function without filling the Id field.
Unless you want to update the Guests field don't call that function again
instead of this, acess the struct field named "Guests"
*/
func (r *Room) FindGuests() GuestMap {
	var db = database.GetConnection()

	query, err := db.Prepare(database.SelectRoomGuests)
	if err != nil {
		panic(err)
	}
	defer query.Close()

	rows, err := query.Query(r.Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if r.Guests == nil {
		r.Guests = make(GuestMap)
	}

	for rows.Next() {
		u := User{}
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Passwd, &u.Img); err != nil {
			panic(err)
		}

		u.ClearCriticalInfo()
		r.Guests[u.Email] = u
	}

	return r.Guests
}

func (r *Room) FindProductsLists() ProductListMap {
	var db = database.GetConnection()

	r.ProductsList = make(ProductListMap)

	rows, err := db.Query(database.SelectProductListByRoom, r.Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		list := ProductList{}

		rows.Scan(&list.Name, &list.Room)

		r.ProductsList[list.Name] = list
	}

	return r.ProductsList
}

func (r *Room) FindTabs() []Tab {
	var db = database.GetConnection()

	r.Tabs = make([]Tab, 0)

	rows, err := db.Query(database.SelectTabsInRoom, r.Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		r.Tabs = append(r.Tabs, Tab{})
		rows.Scan(&r.Tabs[i].Number, &r.Tabs[i].RoomId, &r.Tabs[i].PayValue, &r.Tabs[i].Maded, &r.Tabs[i].Table)
		r.Tabs[i].RemoveMadedTrash()
	}

	return r.Tabs
}

func (r *Room) FindTabsRequests() {
	for i := 0; i < len(r.Tabs); i++ {
		r.Tabs[i].FindRequests()
	}
}

func (r *Room) IsOwner(user User) bool {
	if r.Owner == 0 && r.Id > 0 {
		r.Owner = InitUserByRoom(r.Id).Id
	}

	if r.Owner == user.Id {
		return true
	}

	return false
}

/*
For guarantee that no other guest has been added,
call the FindeGuests() function before using this specific one
*/
func (r *Room) GuestPermission(user User) int {
	if r.Guests == nil {
		r.Guests = make(GuestMap)
	}

	r.FindGuests()

	if !r.IsGuest(user) {
		return 0
	}

	return SelectGuestPermission(user.Id, r.Id)
}

func (r *Room) IsGuest(user User) bool {
	if r.Guests == nil {
		r.Guests = make(GuestMap)
	}

	r.FindGuests()

	if r.Guests[user.Email].Id != 0 {
		return true
	}

	return false
}

func RoomByItsOwner(owner int) Room {
	var room Room
	var db = database.GetConnection()

	rows, err := db.Query(database.SelectRoomByOwner, owner)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&room.Id, &room.Owner)
	}

	return room
}

func RoomByItsId(id int) Room {
	var room Room
	var db = database.GetConnection()

	room.Guests = make(GuestMap)

	rows, err := db.Query(database.SelectRoomById, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&room.Id, &room.Owner)
	}

	return room
}
