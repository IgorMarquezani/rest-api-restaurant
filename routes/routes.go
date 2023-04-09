package routes

import (
	"log"
	"net/http"

	"github.com/api/controllers/invites"
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
	r.HandleFunc("/api/user/register", users.Register).Methods(http.MethodPost)
	r.Handle("/api/user/login", users.Login{}).Methods(http.MethodPost)
	r.Handle("/api/user/auth", session.HandleAuthentication{}).Methods(http.MethodPost)
	r.HandleFunc("/api/user/full-info", users.FullInfo).Methods(http.MethodGet, http.MethodPost)
}

func SessionsRoutes(r *mux.Router) {
	r.Handle("/api/session/auth", session.HandleAuthentication{}).Methods(http.MethodPost)
	r.HandleFunc("/api/session/update/room", session.UpdateActiveRoom).Methods(http.MethodPost, http.MethodPut)
}

func RoomsRoutes(r *mux.Router) {
	r.Handle("/api/room/full-info/", rooms.FullInfo{}).Methods(http.MethodGet)
}

func ProdListRoutes(r *mux.Router) {
	r.HandleFunc("/api/product_list/register", product_list.Register).Methods(http.MethodPost)
}

func ProductsRoutes(r *mux.Router) {
	r.HandleFunc("/api/product/register", products.Register).Methods(http.MethodPost)
	r.Handle("/api/product/update", products.UpdateProduct{}).Methods(http.MethodPut)
	r.HandleFunc("/api/product/{name}/{room}", products.GetProduct).Methods(http.MethodGet)
}

func TabsRoutes(r *mux.Router) {
	r.HandleFunc("/api/tab/register", tabs.Register).Methods(http.MethodPost)
}

func InvitesRoutes(r *mux.Router) {
	r.HandleFunc("/api/invite/send/{email}/", invites.Send).Methods(http.MethodPost)
	r.HandleFunc("/api/invite/accept/{id}/", invites.Send).Methods(http.MethodPost)
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

	log.Fatal(http.ListenAndServe(":3300", r))
}
