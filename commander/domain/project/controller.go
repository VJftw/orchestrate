package project

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"github.com/vjftw/orchestrate/commander/domain/user"
	"github.com/vjftw/orchestrate/commander/middlewares"
)

// Project - Handles actions that can be performed on Projects
type Controller struct {
	render           *render.Render
	UserManager      user.Manager  `inject:"user.manager"`
	UserProvider     user.Provider `inject:"user.provider"`
	ProjectManager   Manager       `inject:"project.manager"`
	ProjectProvider  Provider      `inject:"project.provider"`
	ProjectResolver  Resolver      `inject:"project.resolver"`
	ProjectValidator Validator     `inject:"project.validator"`
}

// Setup - Sets up the controller on the router and a renderer
func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.Handle("/v1/projects", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.securedPostHandler)),
	)).Methods("POST")
}

func (c Controller) securedPostHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUserUUID := context.Get(r, "userUUID")

	// get User via authenticatedUserUUID
	user := c.UserProvider.New()

	err := c.UserManager.GetInto(user, "uuid = ?", authenticatedUserUUID)
	if err != nil {
		c.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	// Unmarshal request into project variable
	project := c.ProjectProvider.New()
	err = c.ProjectResolver.FromRequest(project, r.Body)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// validate the project variable
	res := c.ProjectValidator.Validate(project)
	if res == false {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	project.UUID = uuid.NewV4().String()

	// Save the project variable
	c.ProjectManager.Save(project)

	c.render.JSON(w, http.StatusCreated, project.ToMap())
}
