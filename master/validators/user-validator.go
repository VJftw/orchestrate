package validators

import (
	"github.com/vjftw/orchestrate/master/messages"
	"github.com/vjftw/orchestrate/master/models"
)

type UserValidator struct {
}

// Validate - Validates a User Model
func (uV UserValidator) Validate(userModel models.Model) (bool, *messages.ValidationMessage) {
	return true, nil
}
