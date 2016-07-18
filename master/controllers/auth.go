package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/vjftw/orchestrate/master/managers"
	"github.com/vjftw/orchestrate/master/providers"
	"github.com/vjftw/orchestrate/master/resolvers"
)

// Auth - Handles authentication
type Auth struct {
	render            *render.Render
	UserProvider      providers.IUser      `inject:"provider.user"`
	UserResolver      resolvers.IUser      `inject:"resolver.user"`
	ModelManager      managers.IManager    `inject:"manager.default"`
	AuthTokenProvider providers.IAuthToken `inject:"provider.auth_token"`
}

func (aC *Auth) Setup(router *mux.Router, renderer *render.Render) {
	aC.render = renderer

	router.
		HandleFunc("/v1/auth", aC.authHandler).
		Methods("POST")
}

func (aC Auth) authHandler(w http.ResponseWriter, r *http.Request) {
	user := aC.UserProvider.New()

	// Unmarshal request into user variable
	err := aC.UserResolver.FromRequest(user, r.Body)
	if err != nil {
		aC.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	err = aC.ModelManager.GetInto(user, "email_address = ?", user.EmailAddress)
	if err != nil {
		aC.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	if user.VerifyPassword() {
		token, err := aC.AuthTokenProvider.NewFromUser(user)

		if err != nil {
			aC.render.JSON(w, http.StatusBadRequest, nil)
			return
		}

		aC.render.JSON(w, http.StatusCreated, token.ToMap())
		return
	}

	aC.render.JSON(w, http.StatusUnauthorized, nil)
}
