package models

import (
	"database/sql"
	"fmt"

	"github.com/api/database"
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
	defer insert.Close()
	if err != nil {
		return session, err
	}

	search, err := db.Query(database.SelectSessionByHash, securePS)
	if err != nil {
		return session, err
	}

	if search.Next() {
		search.Scan(&session.Who, &session.ActiveRoom, &session.SecurePS)
	}

	return session, nil
}

func ThereIsSession(u User) (UserSession, bool) {
	var search *sql.Rows
	var err error
	var session UserSession
	var db = database.GetConnection()

	search, err = db.Query(database.SelectSessionByUId, u.Id)
	defer search.Close()
	if err != nil {
		panic(err)
	}

	if search.Next() {
		err := search.Scan(&session.Who, &session.ActiveRoom, &session.SecurePS)
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
		return fmt.Errorf("no session initialized")
	}

	if room.Id <= 0 {
		return fmt.Errorf("room not initialized")
	}

	_, err := db.Query(database.UpdateSessionAcRoom, room.Id, seesion.Who)
	if err != nil {
		return err
	}
	seesion.ActiveRoom = room.Id

	return nil
}

func SessionBySecurePS(hash []byte) (UserSession, bool) {
	var UserSession UserSession
	var db = database.GetConnection()

	search, err := db.Query(database.SelectSessionByHash, hash)
	if err != nil {
		panic(err)
	}

	if search.Next() {
		search.Scan(&UserSession.Who, &UserSession.ActiveRoom, &UserSession.SecurePS)
		return UserSession, true
	}

	return UserSession, false
}
