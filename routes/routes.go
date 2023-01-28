package routes

import (
	"log"
	"net/http"

	"github.com/api/controllers"
	"github.com/api/middleware"
	"github.com/gorilla/mux"
)

func HandleRequest() {
  r := mux.NewRouter()
  r.Use(middleware.SetAllContentType)
  r.HandleFunc("/api/user/register", controllers.NewUser).Methods("Post")
  r.HandleFunc("/api/user/login", controllers.EnterUser).Methods("Post")
  r.HandleFunc("/api/user/update", controllers.UpdateUser).Methods("Put")
  r.HandleFunc("/api/product_list/register", controllers.NewProductList).Methods("Post")
  r.HandleFunc("/api/product/register", controllers.NewProduct).Methods("Post")
  //r.HandleFunc("/api/product/update", controllers.UpdateProduct).Methods("Put")
  
  log.Fatal(http.ListenAndServe(":6000", r))
  
}
