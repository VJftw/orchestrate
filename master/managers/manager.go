package managers

import "github.com/vjftw/orchestrate/master/models"

// Manager interface
type Manager interface {
	Save(models.Model)
	GetInto(models.Model, interface{}, ...interface{}) bool
	Delete(models.Model) bool
}
