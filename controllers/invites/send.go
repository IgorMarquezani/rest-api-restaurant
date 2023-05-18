package invites

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/api/controllers"
	"github.com/api/models"
	"github.com/gorilla/mux"
)

func Send(w http.ResponseWriter, r *http.Request) {
	var (
		room   models.Room
		target models.User
		invite models.Invite
	)

	controllers.AllowCrossOrigin(&w, "*")

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := controllers.ValidJSONFormat(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Unmarshal(body, &invite)

	if invite.InvitingRoom > 0 {
		room = models.RoomByItsId(invite.InvitingRoom)
	} else {
		room = models.RoomByItsId(session.ActiveRoom)
		invite.InvitingRoom = room.Id
	}

	email := mux.Vars(r)["email"]
	_, err = mail.ParseAddress(email)
	if err != nil {
		http.Error(w, "Invalid e-mail format", http.StatusBadRequest)
		return
	}

	if !room.IsOwner(user) {
		if !room.IsGuest(user) {
			http.Error(w, "Not a guest", http.StatusUnauthorized)
			return
		}

		if room.GuestPermission(user) < 2 {
			http.Error(w, "Don't have permission for this operation", http.StatusUnauthorized)
			return
		}
	}

	target = models.SelectUser(email)
	if target.Id == 0 {
		http.Error(w, models.ErrNoSuchUser, http.StatusNotFound)
		return
	}

	if target.Id == room.Owner {
		http.Error(w, "Cannot invite the owner to the room", http.StatusBadRequest)
		return
	}

  if room.IsGuest(target) {
    http.Error(w, "Already a guest", http.StatusAlreadyReported)
    return
  }

	if invite, _ := models.SelectInvite(target, room); invite.Id != 0 {
		http.Error(w, models.ErrInvitedAlready, http.StatusAlreadyReported)
		return
	}

	if invite.Permission > 0 && invite.Permission <= 3 {
		models.InsertInvite(target, room.Id, invite.Permission)
	} else {
		http.Error(w, models.ErrInvalidPermission, http.StatusBadRequest)
		return
	}

	target.ClearCriticalInfo()

	w.WriteHeader(http.StatusOK)
  controllers.EncodeJSON(w, target)
}
