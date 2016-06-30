package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vjftw/orchestrate/master/managers"
	"github.com/vjftw/orchestrate/master/models"
)

// AuthController - Handles authentication
type AuthController struct {
	EntityManager managers.EntityManager `inject:"inline"`
}

// AddRoutes - Adds the routes assosciated to this controller
func (aC AuthController) AddRoutes(r *mux.Router) {
	r.
		HandleFunc("/v1/auth", aC.authHandler).
		Methods("POST")
}

func (aC AuthController) authHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Unmarshal request into user variable
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// 400 on Error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	aC.EntityManager.ORM.FindInto(&user, "email_address = ?", user.EmailAddress)

}
