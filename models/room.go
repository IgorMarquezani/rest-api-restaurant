package models

import (
	"github.com/api/database"
)

type Room struct {
	Id     int `json:"id"`
	Owner  int `json:"owner"`
	Guests guestMap
  ProductsList ProdListMap 
}

// unless you want to update the Guests field don't call that function again
// instead of this, acess the struct field named "Guests"
func (r *Room) FindGuests() guestMap {
	var db = database.GetConnection()

	query, err := db.Prepare(database.SelectRoomGuests)
	if err != nil {
		panic(err)
	}

	result, err := query.Query(r.Id)
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

func SelectRoom(owner int) Room {
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

func RoomById(id int) Room {
  var room Room
  room.Guests = make(guestMap)
  var db = database.GetConnection()

  search, err := db.Query(database.SelectRoomById, id)
  if err != nil {
    panic(err)
  }

  if search.Next() {
    search.Scan(&room.Id, &room.Owner)
  }

  return room
}
