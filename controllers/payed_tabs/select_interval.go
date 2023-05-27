package payedtabs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/api/controllers"
	"github.com/api/models"
	"github.com/gorilla/mux"
)

func SelectInterval(w http.ResponseWriter, r *http.Request) {
	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, models.ErrNoSession, http.StatusUnauthorized)
		return
	}

	room := models.RoomByItsId(session.ActiveRoom)

	ok, permission := room.IsOwnerOrGuest(user)
	if !ok {
		http.Error(w, models.ErrNotAGuest, http.StatusUnauthorized)
		return
	}

	if permission < 2 {
		http.Error(w, models.ErrInsufficientPermission, http.StatusUnauthorized)
		return
	}

	from := mux.Vars(r)["from"]

	dateFrom, err := time.Parse("2006-01-02", from)
	if err != nil {
		http.Error(w, "Date format should be like the following type: yyyy-mm-dd", http.StatusBadRequest)
		return
	}

	to := mux.Vars(r)["to"]

	dateTo, err := time.Parse("2006-01-02", to)
	if err != nil {
		http.Error(w, "Date format should be like the following type: yyyy-mm-dd", http.StatusBadRequest)
		return
	}

	if !dateFrom.Before(dateTo) {
		err := fmt.Sprintf("from %s to %s is a invalid date interval because %s is after %s", from, to, from, to)

		http.Error(w, err, http.StatusBadRequest)
		return
	}

	payedTabs, err := models.SelectPayedTabs(from, to)
	if err != nil {
		panic(err)
	}

	for i, pt := range payedTabs {
		var err error

		payedTabs[i].Date = pt.Date[0:10]

		payedTabs[i].Requests, err = models.SelectPayedRequests(pt.RoomId, pt.Id)
		if err != nil {
			panic(err)
		}
	}

	controllers.EncodeJSON(w, payedTabs)
}
