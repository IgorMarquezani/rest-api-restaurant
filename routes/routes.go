package routes

import (
	"log"
	"net/http"

	"github.com/api/controllers/product_list"
	"github.com/api/controllers/products"
	"github.com/api/controllers/rooms"
	"github.com/api/controllers/session"
	"github.com/api/controllers/tabs"
	"github.com/api/controllers/users"
	"github.com/api/middleware"
	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	r.Handle("/api/user/register", users.Register{}).Methods("POST")
	r.Handle("/api/user/login", users.Login{}).Methods("POST")
	r.Handle("/api/user/auth", session.HandleAuthentication{}).Methods("POST")
	r.Handle("/api/user/full-info", users.UserFullInfo{}).Methods("POST")
}

func SessionsRoutes(r *mux.Router) {
	r.Handle("/api/session/auth", session.HandleAuthentication{}).Methods("POST")
	r.HandleFunc("/api/session/update/room", session.MustUpdateActiveRoom).Methods("POST")
}

func RoomRoutes(r *mux.Router) {
	r.Handle("/api/room/full-info/", rooms.RoomInfo{}).Methods("GET")
}

func ProdListRoutes(r *mux.Router) {
	r.Handle("/api/product_list/register", prodlist.RegisterList{}).Methods("POST")
}

func ProductsRoutes(r *mux.Router) {
	r.Handle("/api/product/register", products.HandleProductRegister{}).Methods("POST")
	r.Handle("/api/product/update", products.UpdateProduct{}).Methods("PUT")
}

func TabsRoutes(r *mux.Router) {
	r.Handle("/api/tab/register", tabs.HandleTabRegister{}).Methods("POST")
}

func HandleRequest() {
	r := mux.NewRouter()
	r.Use(middleware.SetAllContentType)

	UserRoutes(r)
	SessionsRoutes(r)
	RoomRoutes(r)
	ProdListRoutes(r)
	ProductsRoutes(r)
	TabsRoutes(r)

	log.Fatal(http.ListenAndServe(":3300", r))
}
