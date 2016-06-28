package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vjftw/orchestrate/master/managers"
	"github.com/vjftw/orchestrate/master/models"
)

// UserController - Handles actions that can be performed on Users
type UserController struct {
	EntityManager managers.EntityManager `inject:"inline"`
}

// AddRoutes - Adds the routes assosciated to this controller
func (uC *UserController) AddRoutes(r *mux.Router) {
	r.
		HandleFunc("/v1/users", uC.postHandler).
		Methods("POST")
}

func (uC *UserController) postHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Unmarshal request into user variable
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// 400 on Error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the user variable
	resultMap := uC.EntityManager.Validate(&user)
	if resultMap != nil {
		// Respond(w, http.StatusBadRequest, resultMap)
		return
	}

	// Persist the user variable
	uC.EntityManager.Save(&user)

	// write the user variable to output and set http header to 201
	Respond(w, http.StatusCreated, user)
}
