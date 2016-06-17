package managers

import (
	"github.com/vjftw/orchestrate/master/models"
	"github.com/vjftw/orchestrate/master/persisters"
)

// UserManager - Manages the lifecycle of User models
type UserManager struct {
	ORM persisters.IPersister `inject:"persister gorm"`
}

// Save - Persist a new or existing User model. May be stored on multiple storage backends (PGSQL, Redis, etc.)
func (uM *UserManager) Save(user *models.User) {
	uM.ORM.Save(user)
}

// Validate - Validates a given User model.
// Failing Rules:
//   - TODO: Invalid Email Address
//   - TODO: Empty Password
func (uM *UserManager) Validate(user *models.User) error {
	return nil
}
