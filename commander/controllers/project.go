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

// Project - Handles actions that can be performed on Projects
type Project struct {
	render           *render.Render
	ModelManager     managers.IManager   `inject:"manager.default"`
	UserProvider     providers.IUser     `inject:"provider.user"`
	ProjectProvider  providers.IProject  `inject:"provider.project"`
	ProjectResolver  resolvers.IProject  `inject:"resolver.project"`
	ProjectValidator validators.IProject `inject:"validator.project"`
}

// Setup - Sets up the controller on the router and a renderer
func (pC *Project) Setup(router *mux.Router, renderer *render.Render) {
	pC.render = renderer

	router.Handle("/v1/projects", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(pC.securedPostHandler)),
	)).Methods("POST")
}

func (pC Project) securedPostHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUserUUID := context.Get(r, "userUUID")

	// get User via authenticatedUserUUID
	user := pC.UserProvider.New()

	err := pC.ModelManager.GetInto(user, "uuid = ?", authenticatedUserUUID)
	if err != nil {
		pC.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	// Unmarshal request into project variable
	project := pC.ProjectProvider.New()
	err = pC.ProjectResolver.FromRequest(project, r.Body)
	if err != nil {
		pC.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// validate the project variable
	res := pC.ProjectValidator.Validate(project)
	if res == false {
		pC.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	project.UUID = uuid.NewV4().String()

	// Save the project variable
	pC.ModelManager.Save(project)

	pC.render.JSON(w, http.StatusCreated, project.ToMap())
}
