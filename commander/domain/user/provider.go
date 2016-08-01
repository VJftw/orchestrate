package user

import (
	"net/http"

	"github.com/gorilla/context"
)

type Provider interface {
	New() *User
	FromAuthenticatedRequest(*http.Request) (*User, error)
}

type UserProvider struct {
	UserManager *UserManager `inject:"user.manager"`
}

func NewProvider() Provider {
	return &UserProvider{}
}

func (p UserProvider) New() *User {
	return &User{}
}

func (p UserProvider) FromAuthenticatedRequest(r *http.Request) (*User, error) {
	authenticatedUserUUID := context.Get(r, "userUUID")

	user, err := p.UserManager.FindByUUID(authenticatedUserUUID.(string))
	if err != nil {
		return nil, err
	}

	return user, nil
}
