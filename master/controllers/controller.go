package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/vjftw/orchestrate/master/models"
)

// Controller - Interface that defines methods that all controllers should have
type Controller interface {
	AddRoutes(mux.Router)
}

// Respond - Writes the given status code and object to the response
func Respond(w http.ResponseWriter, code int, v models.Serializable) {

	r := render.New(render.Options{
		IndentJSON: true,
	})

	r.JSON(w, code, v.ToMap())
}

// RespondNoBody - Writes the given status code
func RespondNoBody(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func JWTMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	tokenString, err := fromAuthHeader(r)
	if err != nil {
		RespondNoBody(w, http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			RespondNoBody(w, http.StatusUnauthorized)
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("hmacSecret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		context.Set(r, "userID", claims["userId"])
	} else {
		fmt.Println(err)
	}
	next(w, r)
}

// FromAuthHeader is a "TokenExtractor" that takes a give request and extracts
// the JWT token from the Authorization header.
func fromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}
