package tabs

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/api/controllers"
	"github.com/api/models"
	"nhooyr.io/websocket"
)

type tabsChan chan models.Tab

type roomsTabsChan map[int]tabsChan

var websocketsChan roomsTabsChan = make(roomsTabsChan)

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
	if ch, ok := websocketsChan[roomId]; !ok {
		return
	} else {
		ch <- tab
	}
}

func Websocket(w http.ResponseWriter, r *http.Request) {
	var ch tabsChan

	err, _, session := controllers.VerifySession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

  option := strings.Split(r.Header.Get("Origin"), "/")[2]
  options := &websocket.AcceptOptions{OriginPatterns: []string{option}}

  c, err := websocket.Accept(w, r, options)
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusGoingAway, "Closing connection")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = c.Write(ctx, websocket.MessageText, []byte("hello browser"))
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
				c.Write(ctx, websocket.MessageText, []byte(err.Error()))
				continue
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
