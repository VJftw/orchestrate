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
)

// UserController - Handles actions that can be performed on Users
type UserController struct {
	EntityManager managers.EntityManager `inject:"inline"`
	// HashIDService *services.HashIDService `inject:"hashids"`
}

func NewUserController(r *mux.Router) *UserController {
	userController := UserController{}

	r.
		HandleFunc("/v1/users", userController.postHandler).
		Methods("POST")

	r.Handle("/v1/users/{id}", negroni.New(
		negroni.HandlerFunc(JWTMiddleware),
		negroni.Wrap(http.HandlerFunc(userController.putHandlerSec)),
	)).Methods("PUT")

	return &userController
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
	// user.HashID = uC.HashIDService.GenerateHash(int(user.ID))

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

	// get User via userID
	var user models.User
	// uC.EntityManager.ORM.FindInto(user, "hash_id = ?", userID)

	if user.ID > 0 {
		// Unmarshal request into user variable
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			// 400 on Error
			RespondNoBody(w, http.StatusBadRequest)
			return
		}

		if len(user.Password) > 0 {
			user.EncryptPassword()
		}

		uC.EntityManager.Save(&user)

		Respond(w, http.StatusOK, user)
	}

	RespondNoBody(w, http.StatusUnauthorized)
}
