package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/urfave/negroni"
	"github.com/vjftw/orchestrate/master/managers"
	"github.com/vjftw/orchestrate/master/models"
	"github.com/vjftw/orchestrate/master/validators"
)

// UserController - Handles actions that can be performed on Users
type UserController struct {
	ModelManager  managers.Manager     `inject:"manager.default"`
	UserValidator validators.Validator `inject:"validator.user"`
}

// NewUserController - Returns a new UserController
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
	res, vM := uC.UserValidator.Validate(&user)
	if res == false {
		Respond(w, http.StatusBadRequest, vM)
		return
	}

	// Encrypt Password
	user.EncryptPassword()

	user.UUID = uuid.NewV4().Bytes()
	// Persist the user variable
	uC.ModelManager.Save(&user)

	// write the user variable to output and set http header to 201
	Respond(w, http.StatusCreated, user)
}

func (uC UserController) putHandlerSec(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userUUID := vars["id"]
	authenticatedUserUUID := context.Get(r, "userUUID")

	// Quick check for route
	if userUUID != authenticatedUserUUID {
		RespondNoBody(w, http.StatusForbidden)
		return
	}

	// get User via userUUID
	var user models.User
	uC.ModelManager.GetInto(user, "uuid = ?", userUUID)

	if len(user.EmailAddress) > 0 {
		// Unmarshal request into user variable
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			// 400 on Error
			RespondNoBody(w, http.StatusBadRequest)
			return
		}

		// validate the user variable
		res, vM := uC.UserValidator.Validate(&user)
		if res == false {
			Respond(w, http.StatusBadRequest, vM)
			return
		}

		if len(user.Password) > 0 {
			user.EncryptPassword()
		}

		uC.ModelManager.Save(&user)

		Respond(w, http.StatusOK, user)
	}

	RespondNoBody(w, http.StatusUnauthorized)
}
