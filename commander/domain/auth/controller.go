package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/vjftw/orchestrate/commander/domain/user"
)

// Controller - Handles authentication
type Controller struct {
	render       *render.Render
	UserProvider user.Provider `inject:"user.provider"`
	UserResolver user.Resolver `inject:"user.resolver"`
	UserManager  user.Manager  `inject:"user.manager"`
	AuthProvider Provider      `inject:"auth.provider"`
}

func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.
		HandleFunc("/v1/auth", c.authHandler).
		Methods("POST")
}

func (c Controller) authHandler(w http.ResponseWriter, r *http.Request) {
	user := c.UserProvider.New()

	// Unmarshal request into user variable
	err := c.UserResolver.FromRequest(user, r.Body)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	err = c.UserManager.GetInto(user, "email_address = ?", user.EmailAddress)
	if err != nil {
		c.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	if user.VerifyPassword() {
		token := c.AuthProvider.NewFromUser(user)

		// if err != nil {
		// 	c.render.JSON(w, http.StatusBadRequest, nil)
		// 	return
		// }

		c.render.JSON(w, http.StatusCreated, token.ToMap())
		return
	}

	c.render.JSON(w, http.StatusUnauthorized, nil)
}
