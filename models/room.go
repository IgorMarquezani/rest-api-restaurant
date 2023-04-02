package models

import (
	"github.com/api/database"
)

type ProductsInRoom map[int]Products

var RoomProducts ProductsInRoom

type Room struct {
	Id    int `json:"id"`
	Owner int `json:"owner"`

	Guests       GuestMap
	ProductsList ProductListMap
	Tabs         []Tab
}

/*
Don't ever call this function without filling the Id field
unless you want to update the Guests field don't call that function again
instead of this, acess the struct field named "Guests"
*/
func (r *Room) FindGuests() GuestMap {
	var db = database.GetConnection()

	query, err := db.Prepare(database.SelectRoomGuests)
	if err != nil {
		panic(err)
	}
	defer query.Close()

	result, err := query.Query(r.Id)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	if r.Guests == nil {
		r.Guests = make(GuestMap)
	}

	for result.Next() {
		user := User{}
		if err := result.Scan(&user.Id, &user.Name, &user.Email, &user.Passwd, &user.Img); err != nil {
			panic(err)
		}
		user.ClearCriticalInfo()
		r.Guests[user.Email] = user
	}

	return r.Guests
}

func (r *Room) FindProductsLists() ProductListMap {
	var db = database.GetConnection()

	r.ProductsList = make(ProductListMap)

	search, err := db.Query(database.SelectProductListByRoom, r.Id)
	if err != nil {
		panic(err)
	}
  defer search.Close()

	for search.Next() {
		list := ProductList{}
		search.Scan(&list.Name, &list.Room)
		r.ProductsList[list.Name] = list
	}

	return r.ProductsList
}

func (r *Room) FindTabs() []Tab {
	var db = database.GetConnection()
	r.Tabs = make([]Tab, 0)

	search, err := db.Query(database.SelectTabsInRoom, r.Id)
	if err != nil {
		panic(err)
	}
  defer search.Close()

	for i := 0; search.Next(); i++ {
		r.Tabs = append(r.Tabs, Tab{})
		search.Scan(&r.Tabs[i].Number, &r.Tabs[i].RoomId, &r.Tabs[i].PayValue, &r.Tabs[i].Maded, &r.Tabs[i].Table)
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

	result, err := db.Query(database.SelectRoomByOwner, owner)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	if result.Next() {
		result.Scan(&room.Id, &room.Owner)
	}

	return room
}

func RoomByItsId(id int) Room {
	var room Room
	var db = database.GetConnection()
	room.Guests = make(GuestMap)

	search, err := db.Query(database.SelectRoomById, id)
	if err != nil {
		panic(err)
	}
  defer search.Close()

	if search.Next() {
		search.Scan(&room.Id, &room.Owner)
	}

	return room
}
