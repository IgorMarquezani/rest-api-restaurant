package invites

import (
	"net/http"
	"strconv"

	"github.com/api/controllers"
	"github.com/api/models"
	"github.com/gorilla/mux"
)

func Decline(w http.ResponseWriter, r *http.Request) {
	controllers.AllowCrossOrigin(&w, "*")

	err, user, _ := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

	err = models.DeleteInvite(uint(invite.Id))
	if err != nil {
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
	}

  w.WriteHeader(http.StatusOK)
}
