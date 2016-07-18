package validators

import "github.com/vjftw/orchestrate/master/models"

type User struct {
}

// Validate - Validates a User Model
func (uV User) Validate(userModel models.IModel) bool {
	return true
}
