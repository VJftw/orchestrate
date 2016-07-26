package managers

import "github.com/vjftw/orchestrate/commander/models"

// IManager interface
type IManager interface {
	Save(models.IModel) error
	GetInto(models.IModel, interface{}, ...interface{}) error
	Delete(models.IModel) error
}
