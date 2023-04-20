package invites

import (
	"net/http"
	"strconv"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
	"github.com/gorilla/mux"
)

func Accept(w http.ResponseWriter, r *http.Request) {
	err, user, _ := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	controllers.AllowCrossOrigin(&w, "*")

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid invite id", http.StatusNotFound)
		return
	}

	invite, ok := models.SearchForInvites(user.UserInvites(), id)
	if !ok {
		http.Error(w, models.ErrNoSuchInvite, http.StatusNotFound)
		return
	}

	if err := models.InsertGuest(uint(invite.InvitingRoom), uint(invite.Target), invite.Permission); err != nil {
		if database.IsDuplicateKeyError(err.Error()) {
			http.Error(w, "User already in room", http.StatusAlreadyReported)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
