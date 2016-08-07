package cadetGroup

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"github.com/vjftw/orchestrate/commander/domain/project"
	"github.com/vjftw/orchestrate/commander/domain/user"
	"github.com/vjftw/orchestrate/commander/middlewares"
)

// Controller - Handles actions that can be performed on CadetGroups
type Controller struct {
	render              *render.Render
	UserManager         user.Manager    `inject:"user.manager"`
	ProjectManager      project.Manager `inject:"project.manager"`
	CadetGroupResolver  Resolver        `inject:"cadetGroup.resolver"`
	CadetGroupValidator Validator       `inject:"cadetGroup.validator"`
	CadetGroupManager   Manager         `inject:"cadetGroup.manager"`
}

// Setup - Sets up the controller on the router and a renderer
func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.Handle("/v1/projects/{projectUUID}/cadetGroups", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.securedPostHandler)),
	)).Methods("POST")

	router.Handle("/v1/projects/{projectUUID}/cadetGroups", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.securedGetHandler)),
	)).Methods("GET")
}

func (c Controller) securedPostHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUserUUID := context.Get(r, "userUUID")

	user, err := c.UserManager.FindByUUID(authenticatedUserUUID.(string))
	if err != nil {
		c.render.JSON(w, http.StatusUnauthorized, nil)
		return
	}

	projectUUID := mux.Vars(r)["projectUUID"]
	project, err := c.ProjectManager.FindByUserAndUUID(user, projectUUID)
	if err != nil {
		c.render.JSON(w, http.StatusForbidden, nil)
		return
	}

	cadetGroup := c.CadetGroupManager.NewForProject(project)
	err = c.CadetGroupResolver.FromRequest(cadetGroup, r.Body)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	res := c.CadetGroupValidator.Validate(cadetGroup)
	if res == false {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	cadetGroup.UUID = uuid.NewV4().String()
	secureRandom := make([]byte, 10)
	rand.Read(secureRandom)
	keyBytes := sha512.Sum512_256(secureRandom)
	cadetGroup.Key = hex.EncodeToString(keyBytes[:sha512.Size256])

	c.CadetGroupManager.Save(cadetGroup)

	c.render.JSON(w, http.StatusCreated, cadetGroup)
}

func (c Controller) securedGetHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUserUUID := context.Get(r, "userUUID")

	user, err := c.UserManager.FindByUUID(authenticatedUserUUID.(string))
	if err != nil {
		c.render.JSON(w, http.StatusUnauthorized, nil)
		return
	}

	projectUUID := mux.Vars(r)["projectUUID"]

	cadetGroups := c.CadetGroupManager.FindByUserAndProjectUUID(user, projectUUID)

	c.render.JSON(w, http.StatusOK, cadetGroups)
}
