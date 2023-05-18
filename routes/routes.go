package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/api/controllers/invites"
	"github.com/api/controllers/products"
	"github.com/api/controllers/products_lists"
	"github.com/api/controllers/rooms"
	"github.com/api/controllers/sessions"
	"github.com/api/controllers/tabs"
	"github.com/api/controllers/users"
	"github.com/api/middleware"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/api/user/register", users.Register).Methods(http.MethodPost)
	r.Handle("/api/user/login", users.Login{}).Methods(http.MethodPost)
	r.Handle("/api/user/auth", sessions.HandleAuthentication{}).Methods(http.MethodPost)

	r.HandleFunc("/api/user/full-info", users.FullInfo).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/api/user/with-active-room-products", users.WithActiveRoom).Methods(http.MethodGet)

	r.HandleFunc("/api/user/invites", users.Invites).Methods(http.MethodGet)
}

func SessionsRoutes(r *mux.Router) {
	r.Handle("/api/session/auth", sessions.HandleAuthentication{}).Methods(http.MethodPost)
	r.HandleFunc("/api/session/update/room", sessions.UpdateActiveRoom).Methods(http.MethodPost, http.MethodPut)
}

func RoomsRoutes(r *mux.Router) {
	r.Handle("/api/room/full-info/", rooms.FullInfo{}).Methods(http.MethodGet)
}

func ProdListRoutes(r *mux.Router) {
	r.HandleFunc("/api/product_list/register", products_lists.Register).Methods(http.MethodPost)
}

func ProductsRoutes(r *mux.Router) {
	r.HandleFunc("/api/product/register", products.Register).Methods(http.MethodPost)

	r.HandleFunc("/api/product/update/{old-name}", products.Update).Methods(http.MethodPut, http.MethodPost)

	r.HandleFunc("/api/product/select/{name}/{room}", products.GetProduct).Methods(http.MethodGet)
	r.HandleFunc("/api/product/select/{name}", products.GetProduct).Methods(http.MethodGet)

	r.HandleFunc("/api/product/delete/{name}", products.Delete).Methods(http.MethodDelete)
}

func TabsRoutes(r *mux.Router) {
	r.HandleFunc("/api/tab/register", tabs.Register).Methods(http.MethodPost)

	r.HandleFunc("/api/tab/update", tabs.Update).Methods(http.MethodPut)

	r.HandleFunc("/api/tab/delete/{number}/{room}", tabs.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/api/tab/delete/{number}", tabs.Delete).Methods(http.MethodDelete)

	r.HandleFunc("/api/tab/websocket/{room-id}", tabs.Websocket).Methods(http.MethodGet)
}

func InvitesRoutes(r *mux.Router) {
	r.HandleFunc("/api/invite/send/{email}/", invites.Send).Methods(http.MethodPost)

	r.HandleFunc("/api/invite/accept/{id}/", invites.Accept).Methods(http.MethodPost)
}

func HandleRequest() {
	r := mux.NewRouter()
	r.Use(middleware.SetAllContentType)

	SessionsRoutes(r)
	ProdListRoutes(r)
	ProductsRoutes(r)
	InvitesRoutes(r)
	RoomsRoutes(r)
	UserRoutes(r)
	TabsRoutes(r)

  server := http.Server{}
  server.Addr = "localhost:3300"
  server.Handler = r

	fmt.Println("Listening on localhost:3300")
	log.Fatal(server.ListenAndServe())
}
