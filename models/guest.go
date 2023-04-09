package models

import "github.com/api/database"

type GuestMap map[string]User

type Guest struct {
	InvintingRoom   int `json:"invinting_room"`
	UserId          int `json:"user_id"`
	PermissionLevel int `json:"permission_level"`
}

func InsertGuest(roomId, userId, permission uint) error {
  var db = database.GetConnection()

  _, err := db.Query(database.InsertGuest, roomId, userId, permission)
  if err != nil {
    return err
  }

  return nil
}

func SelectGuestPermission(user, room int) int {
	var guest Guest
	var db = database.GetConnection()

	rows, err := db.Query(database.SelectGuestByUserAndRoom, user, room)
	if err != nil {
		panic(err)
	}
  defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&guest.InvintingRoom, &guest.UserId, &guest.PermissionLevel)
		if err != nil {
			panic(err)
		}
	}

	return guest.PermissionLevel
}
