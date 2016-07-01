package persisters

import "github.com/vjftw/orchestrate/master/models"

// IPersister - Persistence functions
type IPersister interface {
	Save(models.IModel)
	FindInto(models.IModel, interface{}, ...interface{})
	Exists(models.IModel, interface{}, ...interface{}) bool
}
