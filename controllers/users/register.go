package users

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/api/database"
	"github.com/api/models"
)

type Register struct {
	userRegister models.UserRegister
}

func (re Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	json.NewDecoder(r.Body).Decode(&re.userRegister)

	_, err := mail.ParseAddress(re.userRegister.User.Email)
	if err != nil {
		http.Error(w, "Invalid E-mail format", http.StatusBadRequest)
		return
	}

	if err := re.userRegister.ThereIsPasswdError(); err != nil {
		http.Error(w, "Invalid passwords", http.StatusBadRequest)
		return
	}

	if err := models.InsertUser(re.userRegister.User); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "E-mail already in use", http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(re.userRegister)
}
