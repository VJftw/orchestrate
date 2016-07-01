package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/vjftw/orchestrate/master/managers"
	"github.com/vjftw/orchestrate/master/models"
	"github.com/vjftw/orchestrate/master/services"
)

// UserController - Handles actions that can be performed on Users
type UserController struct {
	EntityManager managers.EntityManager  `inject:"inline"`
	HashIDService *services.HashIDService `inject:"hashids"`
}

// AddRoutes - Adds the routes assosciated to this controller
func (uC UserController) AddRoutes(r *mux.Router) {
	r.
		HandleFunc("/v1/users", uC.postHandler).
		Methods("POST")

	r.Handle("/v1/users/{id}", negroni.New(
		negroni.HandlerFunc(JWTMiddleware),
		negroni.Wrap(http.HandlerFunc(uC.putHandlerSec)),
	)).Methods("PUT")
}

func (uC UserController) postHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Unmarshal request into user variable
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// 400 on Error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the user variable
	vM := uC.EntityManager.Validate(&user)
	if vM.Valid == false {
		Respond(w, http.StatusBadRequest, vM)
		return
	}

	// Encrypt Password
	user.EncryptPassword()

	// Persist the user variable
	uC.EntityManager.Save(&user)

	// Generate HashID
	user.HashID = uC.HashIDService.GenerateHash(int(user.ID))

	// Persist the user variable
	uC.EntityManager.Save(&user)

	// write the user variable to output and set http header to 201
	Respond(w, http.StatusCreated, user)
}

func (uC UserController) putHandlerSec(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	authenticatedUserID := context.Get(r, "userID")
	fmt.Println(authenticatedUserID)

	// Quick check for route
	if userID != authenticatedUserID {
		RespondNoBody(w, http.StatusForbidden)
		return
	}

	// check if UserID exists in database

	// user := context.Get(r, "user")
	// fmt.Fprintf(w, "This is an authenticated request")
	// fmt.Fprintf(w, "Claim content:\n")
}
