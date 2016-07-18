package providers

import "github.com/vjftw/orchestrate/master/models"

type IUser interface {
	New() *models.User
}

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (uP User) New() *models.User {
	return &models.User{}
}
