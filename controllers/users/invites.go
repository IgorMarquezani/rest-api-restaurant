package users

import (
	"net/http"

	"github.com/api/controllers"
)

func Invites(w http.ResponseWriter, r *http.Request) {
	err, user, _ := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, "Not loged in", http.StatusUnauthorized)
		return
	}

	invites := user.UserInvites()

	w.WriteHeader(http.StatusOK)
	controllers.EncodeJSON(w, invites)
}
