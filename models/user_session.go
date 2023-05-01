package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/api/database"
)

const (
	ErrNoSession string = "No session"
)

type UserSession struct {
	Who        int    `json:"who"`
	ActiveRoom int    `json:"active_room"`
	SecurePS   string `json:"_SecurePS"`
}

func StartSession(u User, securePS string) (UserSession, error) {
	var session UserSession
	var db = database.GetConnection()

	insert, err := db.Query(database.InsertNewSession, u.Id, u.Room.Id, securePS)
	if err != nil {
		return session, err
	}
	defer insert.Close()

	rows, err := db.Query(database.SelectSessionByHash, securePS)
	if err != nil {
		return session, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&session.Who, &session.ActiveRoom, &session.SecurePS)
	}

	return session, nil
}

func ThereIsSession(u User) (UserSession, bool) {
	var db = database.GetConnection()
	var (
		session UserSession
		rows    *sql.Rows
		err     error
	)

	rows, err = db.Query(database.SelectSessionByUId, u.Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&session.Who, &session.ActiveRoom, &session.SecurePS)
		if err != nil {
			panic(err)
		}
		return session, true
	}

	return session, false
}

func UpdateActiveRoom(seesion *UserSession, room Room) error {
	var db = database.GetConnection()

	if _, ok := ThereIsSession(User{Id: seesion.Who}); !ok {
		return fmt.Errorf(ErrNoSession)
	}

	if room.Id <= 0 {
		return fmt.Errorf("Room id value cannot be empty")
	}

	update, err := db.Query(database.UpdateSessionAcRoom, room.Id, seesion.Who)
	if err != nil {
		return err
	}
	update.Close()

	seesion.ActiveRoom = room.Id

	return nil
}

func GetUserActiveRoom(user User) (int, error) {
	var session UserSession
	var db = database.GetConnection()

	row, err := db.Query(database.SelectSessionByUId, user.Id)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	if row.Next() {
		err := row.Scan(&session.Who, &session.ActiveRoom, &session.SecurePS)
		if err != nil {
			panic(err)
		}

		return session.ActiveRoom, err
	}

	return 0, errors.New(ErrNoSession)
}

func SessionBySecurePS(hash []byte) (UserSession, bool) {
	var session UserSession
	var db = database.GetConnection()

	rows, err := db.Query(database.SelectSessionByHash, hash)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&session.Who, &session.ActiveRoom, &session.SecurePS)
		return session, true
	}

	return session, false
}
