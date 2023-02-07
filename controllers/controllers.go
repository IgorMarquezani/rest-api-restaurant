package controllers

import (
	"net/http"

	"github.com/api/models"
	"github.com/api/utils"

)

func verifySessionCookie(r *http.Request) (error, models.User) {
  cookiePS, err := r.Cookie("_SecurePS")
  if err != nil {
    panic(err)
  }

  hash := utils.Invert(cookiePS.Value)

  u, err := models.UserByHash(hash)
  if err != nil {
    return err, u
  }

  return nil, u
}
