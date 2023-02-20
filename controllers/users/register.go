package users

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/api/database"
	"github.com/api/models"
)

type Register struct {
	user models.User
}

func (re Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	json.NewDecoder(r.Body).Decode(&re.user)

  _, err := mail.ParseAddress(re.user.Email)
  if err != nil {
    http.Error(w, "Not valid E-mail format", http.StatusBadRequest)
    return
  }

	if err := models.InsertUser(re.user); err != nil {
    if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "E-mail already in use", http.StatusAlreadyReported)
			return
    }

		http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
	}

	w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(re.user)
}
