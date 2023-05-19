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

	if ok, _ := room.IsOwnerOrGuest(user); !ok {
		http.Error(w, "Not a guest in that room", http.StatusUnauthorized)
		return
	}

	oldTab, err = models.SelectTabByNumber(newTab.Number, room.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if oldTab.Table != newTab.Table {
		if err := models.UpdateTab(oldTab, models.Tab{Table: newTab.Table}); err != nil {
			log.Println(err)
			return
		}
	}

	for i, request := range newTab.Requests {
		if request.Operation == "deleting" {
			product, err := models.SelectOneProduct(room.Id, request.ProductName)
			if err != nil {
				panic(err)
			}

			_, err = models.SelectRequest(request.ProductName, newTab.Number, room.Id)
			if err != nil {
				panic(err)
			}

			models.DecreaseTabValue(newTab.Number, room.Id, product.Price*float64(request.Quantity))

			models.DeleteRequest(newTab.Number, room.Id, request.ProductName)
			continue
		}

		if request.Operation == "updating" {
			request, err := newTab.ToNormalRequest(i)
			if err != nil {
				panic(err)
			}

			oldRequest, err := models.SelectRequest(request.ProductName, newTab.Number, room.Id)
			if err != nil {
				panic(err)
			}

			product, err := models.SelectOneProduct(room.Id, newTab.Requests[i].ProductName)
			if err != nil {
				panic(err)
			}

			if oldRequest.Quantity > request.Quantity {
				models.DecreaseTabValue(newTab.Number, room.Id, product.Price*float64(oldRequest.Quantity-request.Quantity))
			}

			if oldRequest.Quantity < request.Quantity {
				models.IncreaseTabValue(newTab.Number, room.Id, product.Price*float64(request.Quantity-oldRequest.Quantity))
			}

			models.UpdateRequestQuantity(request, uint(request.Quantity))
			continue
		}

		if request.Operation == "inserting" {
			tab := models.Tab{
				Number: newTab.Number,
				RoomId: newTab.RoomId,
			}

			request, err := newTab.ToNormalRequest(i)
			if err != nil {
				panic(err)
			}

			product, err := models.SelectOneProduct(room.Id, newTab.Requests[i].ProductName)
			if err != nil {
				panic(err)
			}

			err = models.InsertRequest(tab, request)
			if err != nil && database.IsDuplicateKeyError(err.Error()) {
				_, err := models.SelectRequest(request.ProductName, newTab.Number, room.Id)
				if err != nil {
					panic(err)
				}

				models.UpdateRequestQuantity(request, uint(request.Quantity))
				continue
			}

			models.IncreaseTabValue(newTab.Number, room.Id, product.Price*float64(request.Quantity))
		}
	}

	tab, _ := models.SelectTabByNumber(oldTab.Number, room.Id)
	SendTab(room.Id, tab)

	w.WriteHeader(http.StatusOK)
}
