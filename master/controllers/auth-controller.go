package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/unrolled/render"
	"github.com/vjftw/orchestrate/master/models"
	"github.com/vjftw/orchestrate/master/models/ephemeral"
	"github.com/vjftw/orchestrate/master/routers"
)

// AuthController - Handles authentication
type AuthController struct {
	render *render.Render
}

func NewAuthController(router *routers.MuxRouter) *AuthController {
	authController := AuthController{
		render: router.Render,
	}

	router.Router.
		HandleFunc("/v1/auth", authController.authHandler).
		Methods("POST")

	return &authController
}

func (aC AuthController) authHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Unmarshal request into user variable
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// 400 on Error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// aC.EntityManager.ORM.FindInto(&user, "email_address = ?", user.EmailAddress)

	// Verify bcrypt Password hash
	if user.VerifyPassword() {
		// if valid, generate JWT with email address and ID in
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": string(user.UUID),
			"nbf":    time.Now().Unix(),
		})

		tokenString, err := token.SignedString([]byte("hmacSecret"))
		if err != nil {
			// 400 on Error
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		authToken := ephemeral.AuthEphemeral{}
		authToken.Token = tokenString

		aC.render.JSON(w, http.StatusCreated, authToken.ToMap())
	}

	// else return 404
	aC.render.JSON(w, http.StatusNotFound, nil)
}
