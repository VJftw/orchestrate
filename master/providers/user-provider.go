package providers

import "github.com/vjftw/orchestrate/master/models"

type IUserProvider interface {
	New() *models.User
}

type UserProvider struct {
}

func NewUserProvider() *UserProvider {
	return &UserProvider{}
}

func (uP UserProvider) New() *models.User {
	return &models.User{}
}
