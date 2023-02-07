package routes

import (
	"log"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/middleware"
	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
}

func ProdListRoutes (r *mux.Router) {

}

func ProductsRoutes (r *mux.Router) {

}

func HandleRequest() {
  r := mux.NewRouter()
  r.Use(middleware.SetAllContentType)
  r.Handle("/api/user/register", controllers.Register{}).Methods("Post")
  r.Handle("/api/user/login", controllers.Login{}).Methods("Post")
  r.Handle("/api/product_list/register", controllers.RegisterList{}).Methods("Post")
  r.HandleFunc("/api/product/register", controllers.NewProduct).Methods("Post")
  r.HandleFunc("/api/product/update", controllers.UpdateProduct).Methods("Put")
  
  log.Fatal(http.ListenAndServe(":6000", r))
}
