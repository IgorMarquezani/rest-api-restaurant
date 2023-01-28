package models

import (

	"github.com/api/database"
)

type Room struct {
  Id int `json:"id"`
  Owner int `json:"owner"` 
  Guests []User
}

// unless you want to update the Guests field don't call that function again
// instead of this, acess the struct field named "Guests"
func (r *Room) FindGuests() []User {
  var guests = make([]User, 30) 
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
    guests = append(guests, user)
  }

  r.Guests = guests

  return guests
}

func SelectRoom(owner int) Room {
  var room Room
  var db = database.GetConnection()

  result, err := db.Query(database.SelectRoomByOwner, owner)
  if err != nil { panic(err) }

  if result.Next() {
    result.Scan(room.Id, room.Owner)
  }

  return room
}
