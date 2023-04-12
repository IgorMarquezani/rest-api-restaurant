package tabs

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/api/controllers"
	"github.com/api/models"
	"nhooyr.io/websocket"
)

var acceptOptions websocket.AcceptOptions = websocket.AcceptOptions{
	OriginPatterns: []string{"localhost:8081"},
}

type tabsChan chan models.Tab

type roomTabsChan map[int]tabsChan

var websocketsChan roomTabsChan = make(roomTabsChan)

var mutex sync.Mutex

func newTabsChan(roomId int) tabsChan {
  var (
    ch tabsChan
    ok bool
  )

	mutex.Lock()
  if ch, ok = websocketsChan[roomId]; !ok {
		websocketsChan[roomId] = make(tabsChan, 20)
    ch = websocketsChan[roomId]
  }

	mutex.Unlock()
	return ch
}

func sendTabInChan(roomId int, tab models.Tab) {
	if _, ok := websocketsChan[roomId]; !ok {
		newTabsChan(roomId)
	}

	websocketsChan[roomId] <- tab
}

func Websocket(w http.ResponseWriter, r *http.Request) {
	var ch tabsChan

	err, _, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	c, err := websocket.Accept(w, r, &acceptOptions)
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusGoingAway, "Closing connection")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = c.Write(ctx, websocket.MessageText, []byte("hello websocket"))
	if err != nil {
		panic(err)
	}

	room := models.RoomByItsId(session.ActiveRoom)

	ch, ok := websocketsChan[room.Id]
	if !ok {
		ch = newTabsChan(room.Id)
	}

	for {
		select {
		case tab := <-ch:
			data, err := json.Marshal(tab)
			if err != nil {
				log.Println("Error: ", err.Error())
			}

			err = c.Write(ctx, websocket.MessageText, data)
			if err != nil {
				log.Println("Error: ", err.Error())
			}
		case <-ctx.Done():
			log.Println("Closing")
			c.Close(websocket.StatusGoingAway, "Closing connection")
		}
	}
}
