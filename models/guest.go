package models

import "github.com/api/database"

type Guest struct {
	InvintingRoom   int `json:"invinting_room"`
	UserId          int `json:"user_id"`
	PermissionLevel int `json:"permission_level"`
}

type GuestMap map[string]User

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
