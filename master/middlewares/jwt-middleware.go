package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/unrolled/render"
)

type JWTMiddleware struct {
	render *render.Render
}

func NewJWTMiddleware(renderer *render.Render) *JWTMiddleware {
	return &JWTMiddleware{
		render: renderer,
	}
}

// JWTMiddleware -
func (m *JWTMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	tokenString, err := fromAuthHeader(r)
	if err != nil {
		m.render.JSON(rw, http.StatusUnauthorized, nil)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, m.render.JSON(rw, http.StatusUnauthorized, nil)
			// return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("hmacSecret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		context.Set(r, "userID", claims["userId"])
	} else {
		fmt.Println(err)
	}
	next(rw, r)
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
