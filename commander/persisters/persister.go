package persisters

import "github.com/vjftw/orchestrate/commander/models"

// IPersister - Persistence functions
type IPersister interface {
	Save(models.IModel) error
	GetInto(models.IModel, interface{}, ...interface{}) error
	Delete(models.IModel) error
}
