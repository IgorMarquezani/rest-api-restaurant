package tabs

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/api/controllers"
	"github.com/api/models"
)

type payer struct {
	Number uint   `json:"number"`
	Date   string `json:"date"`
}

func Pay(w http.ResponseWriter, r *http.Request) {
	var p payer

	controllers.AllowCrossOrigin(&w, "*")

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

	body, err := controllers.ValidJSONFormat(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Unmarshal(body, &p)

	tab, err := models.SelectTabByNumber(int(p.Number), room.Id)
	if err != nil {
		http.Error(w, models.ErrNoSuchTab, http.StatusNoContent)
		return
	}

	_, err = time.Parse("yyyy-mm-dd", p.Date)
	if err != nil {
		http.Error(w, "Date format should be like the following type: yyyy-mm-dd", http.StatusBadRequest)
		return
	}

	pt := models.NewPayedTab(tab, p.Date)

	pt.Id, err = pt.Insert()
	if err != nil {
		panic(err)
	}

	tab.FindRequests()

	pt.AddRequests(tab.Requests)

	if err := pt.InsertRequests(); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	controllers.EncodeJSON(w, pt)
}
