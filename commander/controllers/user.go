package controllers

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"github.com/vjftw/orchestrate/commander/managers"
	"github.com/vjftw/orchestrate/commander/middlewares"
	"github.com/vjftw/orchestrate/commander/providers"
	"github.com/vjftw/orchestrate/commander/resolvers"
	"github.com/vjftw/orchestrate/commander/validators"
)

// User - Handles actions that can be performed on Users
type User struct {
	render        *render.Render
	ModelManager  managers.IManager     `inject:"manager.default"`
	UserValidator validators.IValidator `inject:"validator.user"`
	UserProvider  providers.IUser       `inject:"provider.user"`
	UserResolver  resolvers.IUser       `inject:"resolver.user"`
}

// Setup - Sets up the
func (uC *User) Setup(router *mux.Router, renderer *render.Render) {
	uC.render = renderer

	router.
		HandleFunc("/v1/users", uC.postHandler).
		Methods("POST")

	router.Handle("/v1/users/{id}", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(uC.putHandlerSec)),
	)).Methods("PUT")

}

func (uC User) postHandler(w http.ResponseWriter, r *http.Request) {
	user := uC.UserProvider.New()

	// Unmarshal request into user variable
	err := uC.UserResolver.FromRequest(user, r.Body)
	if err != nil {
		uC.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// validate the user variable
	res := uC.UserValidator.Validate(user)
	if res == false {
		uC.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// Encrypt Password and generate UUID
	user.EncryptPassword()
	user.UUID = uuid.NewV4().String()

	// Persist the user variable
	uC.ModelManager.Save(user)

	// write the user variable to output and set http header to 201
	uC.render.JSON(w, http.StatusCreated, user.ToMap())
}

func (uC User) putHandlerSec(w http.ResponseWriter, r *http.Request) {
	userUUID := mux.Vars(r)["id"]
	authenticatedUserUUID := context.Get(r, "userUUID")

	// Quick check for route
	if userUUID != authenticatedUserUUID {
		uC.render.JSON(w, http.StatusForbidden, nil)
		return
	}

	// get User via userUUID
	user := uC.UserProvider.New()

	err := uC.ModelManager.GetInto(user, "uuid = ?", userUUID)
	if err != nil {
		uC.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	// Unmarshal request into user variable
	err = uC.UserResolver.FromRequest(user, r.Body)
	if err != nil {
		uC.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// validate the user variable
	res := uC.UserValidator.Validate(user)
	if res == false {
		uC.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	if len(user.Password) > 0 {
		user.EncryptPassword()
	}

	uC.ModelManager.Save(user)

	uC.render.JSON(w, http.StatusOK, user.ToMap())

}
