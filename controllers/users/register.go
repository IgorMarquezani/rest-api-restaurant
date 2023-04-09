package users

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/api/database"
	"github.com/api/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var ur models.UserRegister
	json.NewDecoder(r.Body).Decode(&ur)

	_, err := mail.ParseAddress(ur.User.Email)
	if err != nil {
		http.Error(w, "Invalid E-mail format", http.StatusBadRequest)
		return
	}

	if err := ur.ThereIsPasswdError(); err != nil {
		http.Error(w, "Invalid passwords", http.StatusBadRequest)
		return
	}

	if err := models.InsertUser(ur.User); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "E-mail already in use", http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ur.User = models.SelectUser(ur.User.Email)
	ur.User.Room = models.RoomByItsOwner(ur.User.Id)
	if err := models.InsertProductList(models.OrphanList(uint(ur.User.Room.Id))); err != nil {
		panic(err)
	}

	ur.User.ClearCriticalInfo()

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.Encode(ur.User)
}
