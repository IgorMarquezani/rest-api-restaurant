package users

import (
	"encoding/json"
	"github.com/api/models"
	"net/http"
	"strings"
)

type Register struct {
	user models.User
}

func (Re Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	json.NewDecoder(r.Body).Decode(&Re.user)

	if err := models.InsertUser(Re.user); err != nil {
		message := strings.Split(err.Error(), " ")

		if message[1] == "duplicate" || message[2] == "key" {
			http.Error(w, "E-mail already in use", http.StatusAlreadyReported)
			return
		}
		http.Error(w, "Unknow error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
