package tabs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/api/controllers"
	"github.com/api/models"
	"nhooyr.io/websocket"
)

type insertChan chan models.Tab

type deleteChan chan int

type roomClients struct {
	id         int
	insertChan map[string]insertChan
	deleteChan map[string]deleteChan
}

var rooms = make(map[int]*roomClients)

var newRoomClients = make(chan int, 100)

func init() {
	go ListenForNewRooms(newRoomClients)
}

func ListenForNewRooms(ch <-chan int) {
	for {
		select {
		case id := <-ch:
			if _, ok := rooms[id]; !ok {
				rooms[id] = &roomClients{
					id:         id,
					insertChan: make(map[string]insertChan),
					deleteChan: make(map[string]deleteChan),
				}
			}
		}
	}
}

func SendTabToRoom(roomId int, tab models.Tab) {
	if clients, ok := rooms[roomId]; ok {
		for _, ch := range clients.insertChan {
			ch <- tab
		}
	}
}

func DeleteTabInRoom(roomId, tabNumber int) {
	if clients, ok := rooms[roomId]; ok {
		for _, ch := range clients.deleteChan {
			ch <- tabNumber
		}
	}
}

func NewRoom(roomId int, ch chan<- int) *roomClients {
	var (
		room *roomClients
		ok   bool
	)

	ch <- roomId

  // checking loop
  // looping into rooms until it finds the room that the user needs
	for room, ok = rooms[roomId]; !ok; {
		room, ok = rooms[roomId]
	}

	return room
}

func Websocket(w http.ResponseWriter, r *http.Request) {
	var (
		room *roomClients
		ok   bool
	)

	err, user, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if room, ok = rooms[session.ActiveRoom]; !ok {
		room = NewRoom(session.ActiveRoom, newRoomClients)
	}

	if _, ok = room.insertChan[user.Email]; !ok {
		room.insertChan[user.Email] = make(insertChan, 10)
	}

	if _, ok = room.deleteChan[user.Email]; !ok {
		room.deleteChan[user.Email] = make(deleteChan, 10)
	}

	room.HandleClientConnection(w, r, user.Email)
}

func (room *roomClients) HandleClientConnection(w http.ResponseWriter, r *http.Request, userEmail string) {
	option := strings.Split(r.Header.Get("Origin"), "/")[2]
	options := &websocket.AcceptOptions{OriginPatterns: []string{option}}

	conn, err := websocket.Accept(w, r, options)
	if err != nil {
		log.Println("Error during websocket creation:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer conn.Close(websocket.StatusGoingAway, "Closing connection")

	ctx, cancel := context.WithCancel(r.Context())
	ctx = conn.CloseRead(ctx)
	defer cancel()

	insertCh := room.insertChan[userEmail]
	deleteCh := room.deleteChan[userEmail]

	for {
		select {
		case <-ctx.Done():
			log.Println("Closing websocket")
			return

		case tab := <-insertCh:
			data, jsonErr := json.Marshal(tab)
			if jsonErr != nil {
				log.Println("Error:", jsonErr.Error())
				conn.Write(ctx, websocket.MessageType(websocket.StatusInternalError), []byte(jsonErr.Error()))
				continue
			}

			if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
				log.Println("Error:", err.Error())
				return
			}

		case number := <-deleteCh:
			if err := conn.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("Delete tab of number: %d", number))); err != nil {
				log.Println("Error:", err.Error())
				return
			}
		}
	}
}

// Old implementation that only supports one user per room

// type RoomsClientsWebsockets map[int]roomClients

// type roomsTabsChan map[int]tabChan

// var websocketsChan roomsTabsChan = make(roomsTabsChan)

// func newTabsChan(roomId int) tabChan {
// 	var (
// 		ch tabChan
// 		ok bool
// 	)

// 	mutex.Lock()
// 	if ch, ok = websocketsChan[roomId]; !ok {
// 		websocketsChan[roomId] = make(tabChan, 20)
// 		ch = websocketsChan[roomId]
// 	}

// 	mutex.Unlock()
// 	return ch
// }

// func sendTabInChan(roomId int, tab models.Tab) {
// 	if ch, ok := websocketsChan[roomId]; ok {
// 		ch <- tab
// 	}
// }

// func WebsocketOld(w http.ResponseWriter, r *http.Request) {
// 	var (
// 		ch tabChan
// 		ok bool
// 	)

// 	err, _, session := controllers.VerifySession(r)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// 	option := strings.Split(r.Header.Get("Origin"), "/")[2]
// 	options := &websocket.AcceptOptions{OriginPatterns: []string{option}}

// 	conn, err := websocket.Accept(w, r, options)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close(websocket.StatusGoingAway, "Closing connection")

// 	ctx, cancel := context.WithCancel(r.Context())
// 	defer cancel()

// 	ctx = conn.CloseRead(ctx)

// 	room := models.RoomByItsId(session.ActiveRoom)

// 	if ch, ok = websocketsChan[room.Id]; !ok {
// 		ch = newTabsChan(room.Id)
// 	}

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			log.Println("Closing websocket")
// 			return

// 		case tab := <-ch:
// 			data, jsonErr := json.Marshal(tab)
// 			if jsonErr != nil {
// 				log.Println("Error: ", jsonErr.Error())
// 				conn.Write(ctx, websocket.MessageType(websocket.StatusInternalError), []byte(jsonErr.Error()))
// 				continue
// 			}

// 			if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
// 				log.Println("Error: ", err.Error())
// 				return
// 			}
// 		}
// 	}
// }
