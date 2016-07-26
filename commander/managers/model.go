package managers

import (
	"github.com/vjftw/orchestrate/commander/models"
	"github.com/vjftw/orchestrate/commander/persisters"
)

// Model - Default Manager  for models
type Model struct {
	GORMPersister persisters.IPersister `inject:"persister.gorm"`
}

func NewModel() *Model {
	return &Model{}
}

// Save - Saves the model across storages
func (mM Model) Save(m models.IModel) error {
	mM.GORMPersister.Save(m)
	return nil
}

// GetInto - Searches the storages for a model identified by the query and places it into the given model reference.
// Returns true if found, false otherwise
func (mM Model) GetInto(m models.IModel, query interface{}, args ...interface{}) error {
	// check cache

	// check database
	err := mM.GORMPersister.GetInto(m, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Delete - Deletes a model from the storages
func (mM Model) Delete(m models.IModel) error {
	mM.GORMPersister.Delete(m)
	return nil
}
