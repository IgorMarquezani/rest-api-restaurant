package models

import "github.com/api/database"

type UserSession struct {
  Who int `json:"who"`
  SecurePS string `json:"_SecurePS"`
  ActiveRoom int `json:"active_room"`
}

func StartSession(u User, securePS string) error {
  var db = database.GetConnection()

  _, err := db.Query(database.InsertNewSession, u.Id, securePS)

  return err
}
