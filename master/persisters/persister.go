package persisters

import "github.com/vjftw/orchestrate/master/models"

// IPersister - Persistence functions
type IPersister interface {
	Save(models.IModel) error
	GetInto(models.IModel, interface{}, ...interface{}) error
	Delete(models.IModel) error
}
