package validators

import "github.com/vjftw/orchestrate/master/models"

// Validator - Validates a given entity
type Validator interface {
	Validate(models.Model) bool
}
