package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/api/database"
	"github.com/api/utils"
)

const (
	ErrNoSuchUser string = "User does not exist"
)

type UserRegister struct {
	User    User   `json:"user"`
	Confirm string `json:"confirm_passwd"`
}

func (u *UserRegister) ThereIsPasswdError() error {
	if u.User.Passwd == "" || u.Confirm == "" {
		return errors.New("One of the strings is empty")
	}

	if u.User.Passwd != u.Confirm {
		return errors.New("Different passwords")
	}

	return nil
}

type User struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Passwd string `json:"passwd"`
	Img    []byte `json:"img"`

	Room         Room
	ActiveRoom   int
	RoomsAsGuest []Room
	Invites      []Invite
}

func InsertUser(u User) error {
	u.Passwd = utils.HashString(u.Passwd)

	db := database.GetConnection()

	_, err := db.Query(database.InsertUser, u.Name, u.Email, u.Passwd, u.Img)

	return err
}

func (u *User) ClearCriticalInfo() {
	u.Passwd = ""
}

func (u *User) UserInvites() []Invite {
	var db = database.GetConnection()

	rows, err := db.Query(database.SearchInviteByTarget, u.Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	u.Invites = make([]Invite, 0)

	for rows.Next() {
		invite := Invite{}
		err := rows.Scan(&invite.Id, &invite.Target, &invite.InvitingRoom, &invite.Status, &invite.Permission)
		if err != nil {
			panic(err)
		}
		u.Invites = append(u.Invites, invite)
	}

	return u.Invites
}

func (u *User) AcceptedRooms() []Room {
	var db = database.GetConnection()

	rows, err := db.Query(database.SelectGuestRooms, u.Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	u.RoomsAsGuest = make([]Room, 0)
	for rows.Next() {
		room := Room{}
		rows.Scan(&room.Id, &room.Owner)
		u.RoomsAsGuest = append(u.RoomsAsGuest, room)
	}

	return u.RoomsAsGuest
}

func InitUserByRoom(room int) User {
	var u User
	var db = database.GetConnection()

	rows, err := db.Query(database.SelectUserByRoom, room)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Passwd, &u.Img)
		if err != nil {
			panic(err)
		}
	}

	return u
}

// Select User by email or id
func SelectUser(key any) User {
	var (
		err  error
		user User
		rows *sql.Rows
		db   = database.GetConnection()
	)

	switch data := key.(type) {
	case string:
		rows, err = db.Query(database.SearchUserByEmail, data)
	case int:
		rows, err = db.Query(database.SearchUserById, data)
	default:
		panic("Invalid data type")
	}

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Passwd, &user.Img)
		if err != nil {
			panic(err)
		}
	}

	return user
}

func UserBySessionHash(hash string) (User, error) {
	var u User
	db := database.GetConnection()

	rows, err := db.Query(database.SearchUserBySession, hash)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&u.Id, &u.Name, &u.Email, &u.Passwd, &u.Img)
		return u, nil
	}

	return u, fmt.Errorf("No such session")
}
