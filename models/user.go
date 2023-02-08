package models

import (
	"database/sql"
	"fmt"

	"github.com/api/database"
	"github.com/api/utils"
)

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Passwd  string `json:"passwd"`
	Img     []byte `json:"img"`
	Room    Room
	Invites []Invite
	Session UserSession
}

func InitUserByRoom(room int) User {
	var u User
	var db = database.GetConnection()

	query, err := db.Query(database.SelectUserByRoom, room)
	if err != nil {
		panic(err)
	}

	if query.Next() {
		err := query.Scan(&u.Id, &u.Name, &u.Email, &u.Passwd, &u.Img)
		if err != nil {
			panic(err)
		}
	}

	return u
}

func SelectUser(key any) User {
	var search *sql.Rows
	var err error

	db := database.GetConnection()
	u := User{}

	switch data := key.(type) {
	case string:
		search, err = db.Query(database.SearchUserByEmail, data)
	case int:
		search, err = db.Query(database.SearchUserById, data)
	default:
		panic("Invalid data type")
	}

	if err != nil {
		panic(err)
	}

	if search.Next() {
		err = search.Scan(&u.Id, &u.Name, &u.Email, &u.Passwd, &u.Img)
		if err != nil {
			panic(err)
		}
	}
	search.Close()

	return u
}

func InsertUser(u User) error {
	u.Passwd = utils.HashString(u.Passwd)

	db := database.GetConnection()

	_, err := db.Query(database.InsertUser, u.Name, u.Email, u.Passwd, u.Img)

	return err
}

func UserByHash(hash string) (User, error) {
	var u User
	db := database.GetConnection()

	search, err := db.Query(database.SearchUserBySession, hash)
	if err != nil {
		panic(err)
	}

	if search.Next() {
		search.Scan(&u.Id, &u.Name, &u.Email, &u.Passwd, &u.Img)
		u.Session, _ = ThereIsSession(u)
		return u, nil
	}

	return u, fmt.Errorf("No such session with this hash")
}
