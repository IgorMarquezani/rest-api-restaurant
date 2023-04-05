package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func SetAllContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "applicantion/json")
		LogRequest(r)
		next.ServeHTTP(w, r)
	})
}

func LogRequest(r *http.Request) {
	// escape sequence for resiting text color
	var reset string = "\033[0m"
	var urlColor string = "\033[37m"
	var bodyColor string = "\033[35m"
	var methodColor string

	if r.Method == http.MethodPost {
		// Blue
		methodColor = "\033[34m"
	} else if r.Method == http.MethodDelete {
		// Red
		methodColor = "\033[31m"
	} else if r.Method == http.MethodPut {
		// Yellow
		methodColor = "\033[33m"
	} else if r.Method == http.MethodGet {
		// Green
		methodColor = "\033[32m"
	}

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	fmt.Printf("URL: %s%v%s\n", urlColor, r.URL, reset)
	fmt.Printf("Method: %s%v%s\n", methodColor, r.Method, reset)
	fmt.Printf("Body: %s%v%s\n", bodyColor, string(body), reset)
	fmt.Println("+-----------------------------------------------------------+")
}
