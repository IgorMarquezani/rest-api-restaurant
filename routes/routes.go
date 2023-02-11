package routes

import (
	"log"
	"net/http"

	"github.com/api/controllers/product_list"
	"github.com/api/controllers/products"
	"github.com/api/controllers/users"
	"github.com/api/middleware"
	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	r.Handle("/api/user/register", users.Register{}).Methods("POST")
	r.Handle("/api/user/login", users.Login{}).Methods("POST")
}

func ProdListRoutes(r *mux.Router) {
	r.Handle("/api/product_list/register",
		prodlist.RegisterList{}).Methods("POST")
}

func ProductsRoutes(r *mux.Router) {
	r.Handle("/api/product/register",
		products.RegisterProduct{}).Methods("POST")
	r.Handle("/api/product/update",
		products.UpdateProduct{}).Methods("PUT")
}

func HandleRequest() {
	r := mux.NewRouter()
	r.Use(middleware.SetAllContentType)
	UserRoutes(r)
	ProdListRoutes(r)
	ProductsRoutes(r)

	log.Fatal(http.ListenAndServe(":6000", r))
}
