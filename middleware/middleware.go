package middleware

import "net/http"

func SetAllContentType(next http.Handler) http.Handler{
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-type", "applicantion/json")
    next.ServeHTTP(w, r)
  })
}
