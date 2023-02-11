package users

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/api/models"
	"github.com/api/utils"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	user     models.User
	userJson models.User
	session  models.UserSession
}

func (l *Login) setSessionCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    "_SecurePS",
		Value:   utils.Invert(l.session.SecurePS),
		Expires: time.Time{}.Add(time.Minute * 120),
	}

	http.SetCookie(w, &cookie)
}

func (l *Login) newSession() {
	var ok bool
	if l.user.Room.Owner == 0 {
		l.user.Room = models.RoomByItsOwner(l.user.Id)
	}

	if l.session, ok = models.ThereIsSession(l.user); !ok {
		// maybe this should be implemented as a procedure on Postgres
	again:
		var hash string = string(utils.RandomByteArray())
		if _, finded := models.SessionBySecurePS([]byte(hash)); finded {
			goto again
		}

		l.session, _ = models.StartSession(l.user, hash)
		return
	}

	if l.session.ActiveRoom == 0 {
		err := models.UpdateActiveRoom(&l.session, l.user.Room)
		if err != nil {
			panic(err)
		}
	}
}

func (l Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	json.NewDecoder(r.Body).Decode(&l.userJson)
	l.user = models.SelectUser(l.userJson.Email)

	if l.user.Email == "" {
		http.Error(w, "E-mail not registered", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(l.user.Passwd), []byte(l.userJson.Passwd))
	if err != nil {
		http.Error(w, "Incompatible password", http.StatusUnauthorized)
		return
	}

	l.newSession()
	l.setSessionCookie(w)
	l.user.ClearCriticalInfo()

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(l.user)
}
