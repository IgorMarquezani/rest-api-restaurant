package tabs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/api/controllers"
	"github.com/api/models"
	"github.com/gorilla/mux"
	"nhooyr.io/websocket"
)

type insertChan chan models.Tab

type deleteChan chan int

type client struct {
	// client email
	email string
	// to know if the client is listening for new tabs of not
	listening bool
	// chan for new and updated tabs
	inserts insertChan
	// chan for deleted tabs
	deletes deleteChan
}

type room struct {
	// room id
	id int
	// a map of clients using strings as key. Example: clients["someone@gmail.com"]
	clients sync.Map
}

var rooms sync.Map

var mutex sync.Mutex

func SendTab(roomId int, tab models.Tab) {
	var (
		v  any
		r  *room
		ok bool
	)

	if v, ok = rooms.Load(roomId); !ok {
		fmt.Println("No such room of id", roomId)
		return
	}

	if r, ok = v.(*room); !ok {
		fmt.Println("Cannot cast to room pointer", roomId)
		return
	}

	r.clients.Range(func(key, value any) bool {
		var (
			ct *client
			ok bool
		)

		if ct, ok = value.(*client); !ok {
			return false
		}

		if ct.listening {
			ct.inserts <- tab
		}

		return true
	})
}

func DeleteTab(roomId, tabNumber int) {
	var (
		v  any
		r  *room
		ok bool
	)

	if v, ok = rooms.Load(roomId); !ok {
    fmt.Println("No such room of id:", roomId)
		return
	}

	if r, ok = v.(*room); !ok {
		fmt.Println("Cannot cast to room pointer", roomId)
		return
	}

	r.clients.Range(func(key, value any) bool {
		var (
			ct *client
			ok bool
		)

		if ct, ok = value.(*client); !ok {
			return false
		}

		if ct.listening {
			ct.deletes <- tabNumber
		}

		return true
	})
}

func NewRoom(id int) any {
	v, _ := rooms.LoadOrStore(id, &room{
		id:      id,
		clients: sync.Map{},
	})
	fmt.Println("created room of number:", id)

	return v
}

func Websocket(w http.ResponseWriter, rq *http.Request) {
	var (
		v  any
		ok bool
		r  *room
		ct *client
	)

	err, user, _ := controllers.VerifySession(rq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	idStr, _ := mux.Vars(rq)["room-id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid room number", http.StatusBadRequest)
		return
	}

	if room := models.RoomByItsId(id); !room.IsOwner(user) {
		if !room.IsGuest(user) {
			http.Error(w, models.ErrNotAGuest, http.StatusUnauthorized)
			return
		}
	}

	if v, ok = rooms.Load(id); !ok {
		v = NewRoom(id)
	}

	if r, ok = v.(*room); !ok {
		panic("Error getting room")
	}

	if v, ok = r.clients.Load(user.Email); !ok {
		r.clients.LoadOrStore(user.Email, &client{
			email:     user.Email,
			listening: false,
		})
	}

	v, _ = r.clients.Load(user.Email)
	if ct, ok = v.(*client); !ok {
		panic("Error getting the client")
	}

	mutex.Lock()
	if !ct.listening {
		ct.listening = true
		ct.inserts = make(insertChan, 10)
		ct.deletes = make(deleteChan, 10)
	}
	mutex.Unlock()

	r.HandleClientConnection(w, rq, ct)
}

var closeChanMutex sync.Mutex

func closeClientChans(ct *client) {
	closeChanMutex.Lock()
	if ct.listening {
		ct.listening = false
		close(ct.inserts)
		close(ct.deletes)
	}
	closeChanMutex.Unlock()
}

func (room *room) HandleClientConnection(w http.ResponseWriter, r *http.Request, ct *client) {
	option := strings.Split(r.Header.Get("Origin"), "/")[2]
	options := &websocket.AcceptOptions{OriginPatterns: []string{option}}

	conn, err := websocket.Accept(w, r, options)
	if err != nil {
		log.Println("Error during websocket creation:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("New websocket connection")
	defer conn.Close(websocket.StatusGoingAway, "Closing connection")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	ctx = conn.CloseRead(ctx)

outer:
	for {
		select {
		case <-ctx.Done():
			log.Println("Closing websocket")
			conn.Close(websocket.StatusGoingAway, "Closing connection")
			break outer

		case tab := <-ct.inserts:
			log.Println("Tab here from room:", room.id, "Tab:", tab)
			data, jsonErr := json.Marshal(tab)
			if jsonErr != nil {
				log.Println("Error:", jsonErr.Error())
				conn.Write(ctx, websocket.MessageType(websocket.StatusInternalError), []byte(jsonErr.Error()))
				continue
			}

			if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
				log.Println("Error:", err.Error())
				conn.Close(websocket.StatusGoingAway, "Closing connection")
				break outer
			}

		case number := <-ct.deletes:
			log.Println("Delete tav here from room:", room.id)
			if err := conn.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("Delete tab of number: %d", number))); err != nil {
				log.Println("Error:", err.Error())
				conn.Close(websocket.StatusGoingAway, "Closing connection")
				break outer
			}
		}
	}

	closeClientChans(ct)
}
