package providers

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/vjftw/orchestrate/commander/models"
	"github.com/vjftw/orchestrate/commander/models/ephemeral"
)

type IAuthToken interface {
	NewFromUser(user *models.User) (*ephemeral.AuthToken, error)
}

type AuthToken struct {
}

func NewAuthToken() *AuthToken {
	return &AuthToken{}
}

func (a AuthToken) NewFromUser(user *models.User) (*ephemeral.AuthToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userUUID": user.UUID,
		"nbf":      time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte("hmacSecret"))
	if err != nil {
		return nil, err
	}

	return &ephemeral.AuthToken{
		Token: tokenString,
	}, nil
}
