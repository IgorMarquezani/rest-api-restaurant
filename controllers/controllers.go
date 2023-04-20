package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/api/models"
	"github.com/api/utils"
)

const (
	ErrJsonFormat string = "Invalid JSON format"
)

func VerifySessionCookie(r *http.Request) (error, models.User) {
	var user models.User

	cookiePS, err := r.Cookie("_SecurePS")
	if err != nil {
		return err, user
	}

	hash := utils.Invert(cookiePS.Value)

	user, err = models.UserBySessionHash(hash)
	if err != nil {
		return err, user
	}

	return nil, user
}

func VerifySession(r *http.Request) (error, models.User, models.UserSession) {
	var (
		user    models.User
		session models.UserSession
		ok      bool
	)

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

func ValidJSONFormat(reader io.ReadCloser) ([]byte, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if !json.Valid(body) {
		return nil, errors.New(ErrJsonFormat)
	}

	return body, nil
}

func EncodeJSON(w io.Writer, data any) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.Encode(data)
}

func AllowCrossOrigin(w *http.ResponseWriter, origin string) {
	(*w).Header().Set("Access-Control-Allow-Origin", origin)
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin")
}
