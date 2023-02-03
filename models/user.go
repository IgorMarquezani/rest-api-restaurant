package models

import (
	"database/sql"
	"fmt"

	"net/http"

	"github.com/api/database"
)

type User struct {
  Id    int `json:"id"` 
  Name  string `json:"name"`
  Email string `json:"email"`
  Passwd string `json:"passwd"`
  Img   []byte `json:"img"`
  Room Room
  Invites []Invite
}

func InitUserByRoom(room int) User {
  var u User
  var db = database.GetConnection()

  query, err := db.Query(database.SelectUserByRoom, room)
  if err != nil { panic(err) }

  if query.Next() {
    err := query.Scan(&u.Id, &u.Name, &u.Email, &u.Passwd, &u.Img)
    if err != nil { panic(err) }
  }

  return u
}

func InitUserByCookie(r *http.Request) User {
  var u User

  emailCookie, err := r.Cookie("email")
  if err != nil {
    panic(fmt.Errorf("Missing email cookie"))
  }

  u = SelectUser(emailCookie.Value)
  u.Room = SelectRoom(u.Id)
  u.Invites = SelectInviteByUser(u)

  return u
}

func SelectUser(key any) User {
  var search *sql.Rows
  var err     error

  db := database.GetConnection()
  u  := User{}

  switch data := key.(type) {
  case string:
    search, err = db.Query(database.SearchUserByEmail, data)
  case int:
    search, err = db.Query(database.SearchUserById, data)
  default:
    panic("Invalid data type")
  }

  if err != nil { panic(err) }

  if search.Next() {
    err = search.Scan(&u.Id, &u.Name, &u.Email, &u.Passwd, &u.Img)
    if err != nil { panic(err) }
  }

  search.Close()
  return u
}

func InsertUser(u User) error {
  db := database.GetConnection() 
  
  _, err := db.Query(database.InsertUser, u.Name, u.Email, u.Passwd, u.Img)

  return err 
}
