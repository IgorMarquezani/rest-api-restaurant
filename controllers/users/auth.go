package users

import (
	"encoding/json"
	"net/http"

	"github.com/api/controllers"
)

type Authenticate struct {
}

func (au Authenticate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  err, user, _ := controllers.VerifySession(r)
  if err != nil {
    http.Error(w, "Not loged in", http.StatusFailedDependency)
    return
  }

  user.ClearCriticalInfo()

  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(user)
}
