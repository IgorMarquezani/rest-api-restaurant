package models

import "github.com/api/database"

type guestMap map[string]User

type Guest struct {
	InvintingRoom   int `json:"invinting_room"`
	UserId          int `json:"user_id"`
	PermissionLevel int `json:"permission_level"`
}

func SelectGuestPermission(user, room int) int {
  var guest Guest
  var db = database.GetConnection()

  search, err := db.Query(database.SelectGuestByUserAndRoom, user, room)
  if err != nil {
    panic(err)
  }

  if search.Next() {
    err := search.Scan(&guest.InvintingRoom, &guest.UserId, &guest.PermissionLevel)
    if err != nil {
      panic(err)
    }
  }

  return guest.PermissionLevel
}
