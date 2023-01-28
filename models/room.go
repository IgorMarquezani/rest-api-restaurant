package models

import (
	"net/http"
  "strconv"
  "fmt"

	"github.com/api/database"
)

type Room struct {
  Id int `json:"id"`
  Owner int `json:"owner"` 
  OwnerInfo User
  Guests []User
}

// unless you want to update the Guests field don't call that function again
// instead of this, acess the struct field named "Guests"
func (r *Room) FindGuests() []User {
  var guests []User
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

func (r *Room) GetOwnerInfo() User {
  var user User
  var db = database.GetConnection()

  result, err := db.Query(database.SearchUserById, r.Owner)
  if err != nil {
    panic(err)
  }

  if result.Next() {
    if err := result.Scan(&user.Id, &user.Name, &user.Email, &user.Passwd, &user.Img); err != nil {
      panic(err)
    }
  }

  r.OwnerInfo = user

  return user
}

func (r *Room) FindOwnerId() {
  if r.Id == 0 {
    panic("Room has nil id")
  }

  var db = database.GetConnection()

  query, err := db.Query(database.SearchRoomById, r.Id)
  if err != nil {
    panic(err)
  }

  if query.Next() {
    if err := query.Scan(&r.Owner); err != nil {
      panic(err)
    }
  }
}

// not shure functions
func GetRoom(cookie http.Cookie) (int, error) {
  roomStr := cookie.Value
  room, err := strconv.Atoi(roomStr)

  if err != nil {
    return room, fmt.Errorf("The cookie is missing the room value")
  }

  return room, nil
}

func (r *Room) SetValues(cookie http.Cookie) error {
  db := database.GetConnection()

  room, err := GetRoom(cookie)
  if err != nil {
    return err
  }

  search, err := db.Query(database.SelectRoom, room)
  if search.Next() {
    search.Scan(&r.Id, &r.Owner)
  }

  return nil
}
