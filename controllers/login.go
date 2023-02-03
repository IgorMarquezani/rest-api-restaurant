package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/api/database"
	"github.com/api/models"

	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
  u models.User
  db *sql.DB 
}

func (l Login) hashPasswd() string {
  if l.u.Passwd == "" { panic("Passwd is empty") }

  hashed, err := bcrypt.GenerateFromPassword([]byte(l.u.Passwd), 8)
  if err != nil { panic(err) }

  return string(hashed)
}

func (l Login) setPasswdCookie(w http.ResponseWriter, hashed string) {
  err := bcrypt.CompareHashAndPassword( []byte(hashed), []byte(l.u.Passwd))
  if err != nil {
    panic(err)
  }

  cookie := http.Cookie {
    Name:  "_SecurePS",
    Value: hashed,
  }

  http.SetCookie(w, &cookie)
}

func (l Login) setActiveRoom() {
  var session models.UserSession
  l.u.Room = models.SelectRoom(l.u.Id)

  search, err := l.db.Query(database.SelectSessionByUId, l.u.Id)
  if err != nil { panic(err) }

  if search.Next() {
    err := search.Scan(session.Who, session.SecurePS, session.ActiveRoom)
    if err != nil { panic(err) }
  }
  search.Close()

  if session.ActiveRoom != 0 { return }

  _, err = l.db.Query(database.UpdateSessionAcRoom, l.u.Room.Id)
  if err != nil {
    panic(err)
  }
}

func (l Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  var uj models.User
  json.NewDecoder(r.Body).Decode(&uj)
  us := models.SelectUser(uj.Email);

  if us.Email == "" {
    http.Error(w, "E-mail not registered", http.StatusNotFound)
    return
  }

  if uj.Passwd != us.Passwd {
    http.Error(w, "Incompatible password", http.StatusConflict)
    return
  }

  l.u = us
  l.u.Room = models.SelectRoom(l.u.Id)
  l.db = database.GetConnection()
  //l.setActiveRoom()
  l.setPasswdCookie(w, l.hashPasswd())

  w.WriteHeader(http.StatusAccepted)
  json.NewEncoder(w).Encode(us)
}
