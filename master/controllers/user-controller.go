package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"github.com/vjftw/orchestrate/master/managers"
	"github.com/vjftw/orchestrate/master/middlewares"
	"github.com/vjftw/orchestrate/master/models"
	"github.com/vjftw/orchestrate/master/providers"
	"github.com/vjftw/orchestrate/master/routers"
	"github.com/vjftw/orchestrate/master/validators"
)

// UserController - Handles actions that can be performed on Users
type UserController struct {
	render        *render.Render
	ModelManager  managers.Manager        `inject:"manager.default"`
	UserValidator validators.Validator    `inject:"validator.user"`
	UserProvider  providers.IUserProvider `inject:"provider.user"`
}

// NewUserController - Returns a new UserController
func NewUserController(router *routers.MuxRouter) *UserController {
	userController := UserController{
		render: router.Render,
	}

	router.Router.
		HandleFunc("/v1/users", userController.postHandler).
		Methods("POST")

	router.Router.Handle("/v1/users/{id}", negroni.New(
		middlewares.NewJWTMiddleware(router.Render),
		negroni.Wrap(http.HandlerFunc(userController.putHandlerSec)),
	)).Methods("PUT")

	return &userController
}

func (uC UserController) postHandler(w http.ResponseWriter, r *http.Request) {
	user := uC.UserProvider.New()

	// Unmarshal request into user variable
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uC.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// validate the user variable
	res, _ := uC.UserValidator.Validate(user)
	if res == false {
		uC.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// Encrypt Password
	user.EncryptPassword()
	//
	user.UUID = uuid.NewV4().String()
	// Persist the user variable
	uC.ModelManager.Save(user)

	// write the user variable to output and set http header to 201
	uC.render.JSON(w, http.StatusCreated, user.ToMap())
}

func (uC UserController) putHandlerSec(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userUUID := vars["id"]
	authenticatedUserUUID := context.Get(r, "userUUID")

	// Quick check for route
	if userUUID != authenticatedUserUUID {
		uC.render.JSON(w, http.StatusForbidden, nil)
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
			uC.render.JSON(w, http.StatusBadRequest, nil)
			return
		}

		// validate the user variable
		res, vM := uC.UserValidator.Validate(&user)
		if res == false {
			uC.render.JSON(w, http.StatusBadRequest, vM.ToMap())
			return
		}

		if len(user.Password) > 0 {
			user.EncryptPassword()
		}

		uC.ModelManager.Save(&user)

		uC.render.JSON(w, http.StatusOK, user.ToMap())
	}

	uC.render.JSON(w, http.StatusUnauthorized, nil)
}
