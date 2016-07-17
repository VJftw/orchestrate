package controllers

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"github.com/vjftw/orchestrate/master/managers"
	"github.com/vjftw/orchestrate/master/middlewares"
	"github.com/vjftw/orchestrate/master/providers"
	"github.com/vjftw/orchestrate/master/resolvers"
	"github.com/vjftw/orchestrate/master/validators"
)

// UserController - Handles actions that can be performed on Users
type UserController struct {
	render        *render.Render
	ModelManager  managers.Manager        `inject:"manager.default"`
	UserValidator validators.Validator    `inject:"validator.user"`
	UserProvider  providers.IUserProvider `inject:"provider.user"`
	UserResolver  resolvers.IUserResolver `inject:"resolver.user"`
}

// Setup - Sets up the controller
func (uC *UserController) Setup(router *mux.Router, renderer *render.Render) {
	uC.render = renderer

	router.
		HandleFunc("/v1/users", uC.postHandler).
		Methods("POST")

	router.Handle("/v1/users/{id}", negroni.New(
		middlewares.NewJWTMiddleware(renderer),
		negroni.Wrap(http.HandlerFunc(uC.putHandlerSec)),
	)).Methods("PUT")

}

func (uC UserController) postHandler(w http.ResponseWriter, r *http.Request) {
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

func (uC UserController) putHandlerSec(w http.ResponseWriter, r *http.Request) {
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
