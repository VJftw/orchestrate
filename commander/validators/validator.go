package validators

import "github.com/vjftw/orchestrate/commander/models"

// IValidator - Validates a given entity
type IValidator interface {
	Validate(models.IModel) bool
}
