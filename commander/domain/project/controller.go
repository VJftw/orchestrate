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

// Controller - Handles actions that can be performed on Projects
type Controller struct {
	render           *render.Render
	UserManager      user.Manager `inject:"user.manager"`
	ProjectManager   Manager      `inject:"project.manager"`
	ProjectResolver  Resolver     `inject:"project.resolver"`
	ProjectValidator Validator    `inject:"project.validator"`
}

// Setup - Sets up the controller on the router and a renderer
func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.Handle("/v1/projects", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.securedPostHandler)),
	)).Methods("POST")

	router.Handle("/v1/projects", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.securedGetHandler)),
	)).Methods("GET")

	router.Handle("/v1/projects/{projectUUID}", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.securedGetOneHandler)),
	)).Methods("GET")
}

func (c Controller) securedPostHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUserUUID := context.Get(r, "userUUID")

	user, err := c.UserManager.FindByUUID(authenticatedUserUUID.(string))
	if err != nil {
		c.render.JSON(w, http.StatusUnauthorized, nil)
		return
	}

	// Unmarshal request into project variable
	project := c.ProjectManager.NewForUser(user)
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

	c.render.JSON(w, http.StatusCreated, project)
}

func (c Controller) securedGetHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUserUUID := context.Get(r, "userUUID")

	user, err := c.UserManager.FindByUUID(authenticatedUserUUID.(string))
	if err != nil {
		c.render.JSON(w, http.StatusUnauthorized, nil)
		return
	}

	projects := c.ProjectManager.FindByUser(user)

	c.render.JSON(w, http.StatusOK, projects)
}

func (c Controller) securedGetOneHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUserUUID := context.Get(r, "userUUID")

	user, err := c.UserManager.FindByUUID(authenticatedUserUUID.(string))
	if err != nil {
		c.render.JSON(w, http.StatusUnauthorized, nil)
		return
	}

	projectUUID := mux.Vars(r)["projectUUID"]
	project, err := c.ProjectManager.FindByUserAndUUID(user, projectUUID)

	if err != nil {
		c.render.JSON(w, http.StatusNotFound, nil)
	}

	c.render.JSON(w, http.StatusOK, project)
}
