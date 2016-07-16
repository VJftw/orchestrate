package persisters

import "github.com/vjftw/orchestrate/master/models"

// Persister - Persistence functions
type Persister interface {
	Save(models.Model)
	FindInto(models.Model, interface{}, ...interface{})
	Exists(models.Model, interface{}, ...interface{}) bool
	Delete(models.Model) bool
}
