package managers

import (
	"github.com/vjftw/orchestrate/master/models"
	"github.com/vjftw/orchestrate/master/persisters"
)

// EntityManager - Manages the lifecycle of entities
type EntityManager struct {
	ORM persisters.IPersister `inject:"persister gorm"`
}

// Save - Persist a new or existing Entity. May be stored on multiple storage backends (PGSQL, Redis, etc.)
func (eM EntityManager) Save(entity models.IModel) {
	eM.ORM.Save(entity)
}

// Validate - Validates a given Entity.
func (eM EntityManager) Validate(entity models.IModel) error {
	return nil
}
