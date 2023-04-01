package middleware

import (
	"fmt"
	"net/http"

	"github.com/api/controllers"
)

func SetAllContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "applicantion/json")
    controllers.AllowCrossOrigin(&w, "*")
		LogRequest(w, r)
		next.ServeHTTP(w, r)
	})
}

func LogRequest(w http.ResponseWriter, r *http.Request) {
	var color string
	var reset string = "\033[0m"

	if r.Method == http.MethodPost {
		color = "\033[34m"
	}

	if r.Method == http.MethodDelete {
		color = "\033[31m"
	}

	if r.Method == http.MethodPut {
		color = "\033[33m"
	}

	if r.Method == http.MethodGet {
		color = "\033[32m"
	}

	fmt.Printf("%s%v%s URL: %v\n", color, r.Method, reset, r.URL)
	fmt.Println(r.Body)
}
