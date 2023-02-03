package controllers

import (
  "net/http"
  "encoding/json"
  "strings"
  "github.com/api/models"
)

type Register struct {
  u models.User
}

func (Re Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  json.NewDecoder(r.Body).Decode(&Re.u)

  if err := models.InsertUser(Re.u); err != nil {
    message := strings.Split(err.Error(), "")

    if message[1] == "duplicate" || message[2] == "key" {
      http.Error(w, "E-mail already in use", http.StatusAlreadyReported)
      return
    }

    http.Error(w, "Unknow error", http.StatusConflict)
  }
  
  w.WriteHeader(http.StatusCreated)
}
