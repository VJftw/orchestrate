package user

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"github.com/vjftw/orchestrate/commander/middlewares"
)

type Controller struct {
	render        *render.Render
	UserManager   Manager   `inject:"user.manager"`
	UserValidator Validator `inject:"user.validator"`
	UserResolver  Resolver  `inject:"user.resolver"`
}

func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.
		HandleFunc("/v1/users", c.postHandler).
		Methods("POST")

	router.Handle("/v1/users/{id}", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.securePutHandler)),
	)).Methods("PUT")

}

func (c Controller) postHandler(w http.ResponseWriter, r *http.Request) {
	user := c.UserManager.New()

	// Unmarshal request into user variable
	err := c.UserResolver.FromRequest(user, r.Body)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// validate the user variable
	res := c.UserValidator.Validate(user)
	if res == false {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// Encrypt Password and generate UUID
	user.EncryptPassword()
	user.UUID = uuid.NewV4().String()

	// Persist the user variable
	c.UserManager.Save(user)

	// write the user variable to output and set http header to 201
	c.render.JSON(w, http.StatusCreated, user)
}

func (c Controller) securePutHandler(w http.ResponseWriter, r *http.Request) {
	userUUID := mux.Vars(r)["id"]
	authenticatedUserUUID := context.Get(r, "userUUID")

	// Quick check for route
	if userUUID != authenticatedUserUUID {
		c.render.JSON(w, http.StatusForbidden, nil)
		return
	}

	// get User via userUUID
	user, err := c.UserManager.FindByUUID(userUUID)
	if err != nil {
		c.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	// Unmarshal request into user variable
	err = c.UserResolver.FromRequest(user, r.Body)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// validate the user variable
	res := c.UserValidator.Validate(user)
	if res == false {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	if len(user.Password) > 0 {
		user.EncryptPassword()
	}

	c.UserManager.Save(user)

	c.render.JSON(w, http.StatusOK, user)
}
