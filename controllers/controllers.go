package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/api/models"
	"github.com/api/utils"
)

func VerifySessionCookie(r *http.Request) (error, models.User) {
	var user models.User

	cookiePS, err := r.Cookie("_SecurePS")
	if err != nil {
		return err, user
	}

	hash := utils.Invert(cookiePS.Value)
	fmt.Println(hash, "TERMINACAO")

	user, err = models.UserBySessionHash(hash)
	if err != nil {
		return err, user
	}

	return nil, user
}

func VerifySession(r *http.Request) (error, models.User, models.UserSession) {
	var user models.User
	var session models.UserSession
	var ok bool

	err, user := VerifySessionCookie(r)
	if err != nil {
		return err, user, session
	}

	session, ok = models.ThereIsSession(user)
	if !ok {
		return errors.New("Please Log in"), user, session
	}

	return nil, user, session
}

func AllowCrossOrigin(w *http.ResponseWriter, origin string) {
  (*w).Header().Set("Access-Control-Allow-Origin", origin)
  (*w).Header().Set("Access-Control-Allow-Headers", "Origin")
}
