package validators

import "github.com/vjftw/orchestrate/commander/models"

type IUser interface {
	Validate(*models.User) bool
}

type User struct {
}

func NewUser() *User {
	return &User{}
}

// Validate - Validates a User Model
func (uV User) Validate(u *models.User) bool {
	return true
}
