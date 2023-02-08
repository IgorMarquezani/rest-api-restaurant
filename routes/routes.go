package routes

import (
	"log"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/middleware"
	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	r.Handle("/api/user/register", controllers.Register{}).Methods("Post")
	r.Handle("/api/user/login", controllers.Login{}).Methods("Post")
}

func ProdListRoutes(r *mux.Router) {
	r.Handle("/api/product_list/register", controllers.RegisterList{}).Methods("Post")
}

func ProductsRoutes(r *mux.Router) {
	r.Handle("/api/product/register", controllers.RegisterProduct{}).Methods("Post")
}

func HandleRequest() {
	r := mux.NewRouter()
	r.Use(middleware.SetAllContentType)
	UserRoutes(r)
	ProdListRoutes(r)
	ProductsRoutes(r)

	log.Fatal(http.ListenAndServe(":6000", r))
}
