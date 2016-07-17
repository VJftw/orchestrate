package persisters

import "github.com/vjftw/orchestrate/master/models"

// Persister - Persistence functions
type Persister interface {
	Save(models.Model) error
	GetInto(models.Model, interface{}, ...interface{}) error
	Delete(models.Model) error
}
