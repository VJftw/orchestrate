package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// UserController - Handles actions that can be performed on Users
type UserController struct {
}

// AddRoutes - Adds the routes assosciated to this controller
func (uC *UserController) AddRoutes(r *mux.Router) {
	r.
		HandleFunc("/v1/users", uC.postHandler).
		Methods("POST")
}

func (uC *UserController) postHandler(w http.ResponseWriter, r *http.Request) {
	// var user models.User

	// deJSON request into user variable
	// 400 on Error

	// validate the user variable
	// 400 on Error

	// Persist the user variable

	// write the user variable to output and set http header to 201
}
