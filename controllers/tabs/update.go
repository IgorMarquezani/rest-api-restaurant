package tabs

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/database"
	"github.com/api/models"
)

func Update(w http.ResponseWriter, r *http.Request) {
	var (
		newTab models.UpdatingTab
		oldTab models.Tab
		room   models.Room
	)

	controllers.AllowCrossOrigin(&w, r.Header.Get("Origin"))

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := controllers.ValidJSONFormat(r.Body)
	if err != nil {
		http.Error(w, controllers.ErrJsonFormat, http.StatusBadRequest)
		return
	}
	r.Body.Close()

	json.Unmarshal(body, &newTab)
	if newTab.Number <= 0 {
		http.Error(w, models.ErrInvalidTabNumber, http.StatusBadRequest)
		return
	}

	if newTab.RoomId <= 0 {
		room = models.RoomByItsId(session.ActiveRoom)
		newTab.RoomId = room.Id
	} else {
		room = models.RoomByItsId(newTab.RoomId)
	}

	if !room.IsOwner(user) {
		if !room.IsGuest(user) {
			http.Error(w, "Not a guest in that room", http.StatusUnauthorized)
			return
		}
	}

	oldTab, err = models.SelectTabByNumber(newTab.Number, room.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := models.UpdateTab(oldTab, models.Tab{Table: newTab.Table}); err != nil {
		log.Println(err)
		return
	}

	for i := range newTab.Requests {
		if newTab.Requests[i].Operation == "deleting" {
			models.DeleteRequest(newTab.Number, newTab.RoomId, newTab.Requests[i].ProductName)
			continue
		}

		if newTab.Requests[i].Operation == "updating" {
			request := models.Request{
				ProductName: newTab.Requests[i].ProductName,
				TabNumber:   newTab.Number,
				TabRoom:     newTab.RoomId,
			}

			models.UpdateRequestQuantity(request, uint(newTab.Requests[i].Quantity))
			continue
		}

		if newTab.Requests[i].Operation == "inserting" {
			tab := models.Tab{
				Number: newTab.Number,
				RoomId: newTab.RoomId,
			}

			request := models.Request{
				ProductListRoom: newTab.RoomId,
				ProductName:     newTab.Requests[i].ProductName,
				Quantity:        newTab.Requests[i].Quantity,
			}

			err := models.InsertRequest(tab, request)
			if err != nil && database.IsDuplicateKeyError(err.Error()) {
				request := models.SelectRequest(newTab.Requests[i].ProductName, newTab.Number, newTab.RoomId)
				models.UpdateRequestQuantity(request, uint(request.Quantity+newTab.Requests[i].Quantity))
			}
		}
	}

	SendTab(room.Id, newTab.ToNormalTab())

	w.WriteHeader(http.StatusOK)
}
