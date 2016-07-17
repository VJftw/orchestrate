package managers

import (
	"github.com/vjftw/orchestrate/master/models"
	"github.com/vjftw/orchestrate/master/persisters"
)

// ModelManager - Default Manager for models
type ModelManager struct {
	GORMPersister persisters.Persister `inject:"persister.gorm"`
}

// Save - Saves the model across storages
func (mM ModelManager) Save(m models.Model) error {
	mM.GORMPersister.Save(m)
	return nil
}

// GetInto - Searches the storages for a model identified by the query and places it into the given model reference.
// Returns true if found, false otherwise
func (mM ModelManager) GetInto(m models.Model, query interface{}, args ...interface{}) error {
	// check cache

	// check database
	err := mM.GORMPersister.GetInto(m, query, args)
	if err != nil {
		return err
	}

	return nil
}

// Delete - Deletes a model from the storages
func (mM ModelManager) Delete(m models.Model) error {
	mM.GORMPersister.Delete(m)
	return nil
}
