package models

import (
	"net/http"
)

type ProductList struct {
  Name string `json:"name"`
  Room int `json:"room"`
}

func (p *ProductList) SetRoom(cookie http.Cookie) error {
  room, err := GetRoom(cookie)
  if err != nil {
    return err 
  }

  p.Room = room

  return nil
}
