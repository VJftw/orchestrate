package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/vjftw/orchestrate/master/ephemeral"
	"github.com/vjftw/orchestrate/master/managers"
	"github.com/vjftw/orchestrate/master/models"
)

// AuthController - Handles authentication
type AuthController struct {
	EntityManager managers.EntityManager `inject:"inline"`
}

// AddRoutes - Adds the routes assosciated to this controller
func (aC AuthController) AddRoutes(r *mux.Router) {
	r.
		HandleFunc("/v1/auth", aC.authHandler).
		Methods("POST")
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

		Respond(w, http.StatusCreated, authToken)
	}

	// else return 404
	RespondNoBody(w, http.StatusNotFound)
}
