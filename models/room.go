package models

import (
	"github.com/api/database"
)

type Room struct {
	Id            int `json:"id"`
	Owner         int `json:"owner"`

	Guests        GuestMap
	ProductsList ProductListMap
  Tabs []Tab
}

/*
Don't ever call this function without filling the Id field
unless you want to update the Guests field don't call that function again
instead of this, acess the struct field named "Guests"
*/
func (r *Room) FindGuests() GuestMap {
	var db = database.GetConnection()

	query, err := db.Prepare(database.SelectRoomGuests)
	defer query.Close()
	if err != nil {
		panic(err)
	}

	result, err := query.Query(r.Id)
	defer result.Close()
	if err != nil {
		panic(err)
	}

	for result.Next() {
		user := User{}
		if err := result.Scan(&user.Id, &user.Name, &user.Email, &user.Passwd, &user.Img); err != nil {
			panic(err)
		}

		r.Guests[user.Email] = user
	}

	return r.Guests
}

func (r *Room) FindTabs() {
  var db = database.GetConnection()

  search, err := db.Query(database.SelectTabsInRoom, r.Id)
  if err != nil {
    panic(err)
  }

  for i := 0; search.Next(); i++ {
    r.Tabs = append(r.Tabs, Tab{})
    search.Scan(&r.Tabs[i].Number, &r.Tabs[i].RoomId, &r.Tabs[i].PayValue, &r.Tabs[i].Maded)
  }
}

func (r *Room) IsOwner(user User) bool {
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
		r.FindGuests()
	}

	if !r.IsGuest(user) {
		return 0
	}

	return SelectGuestPermission(user.Id, r.Id)
}

func (r *Room) IsGuest(user User) bool {
	if r.Guests == nil {
		r.FindGuests()
	}

	if r.Guests[user.Email].Id != 0 {
		return true
	}

	return false
}

func RoomByItsOwner(owner int) Room {
	var room Room
	var db = database.GetConnection()

	result, err := db.Query(database.SelectRoomByOwner, owner)
	defer result.Close()
	if err != nil {
		panic(err)
	}

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

	if search.Next() {
		search.Scan(&room.Id, &room.Owner)
	}

	return room
}
