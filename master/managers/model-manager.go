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
func (mM ModelManager) Save(m models.Model) {
	mM.GORMPersister.Save(m)
}

// GetInto - Searches the storages for a model identified by the query and places it into the given model reference.
// Returns true if found, false otherwise
func (mM ModelManager) GetInto(m models.Model, query interface{}, args ...interface{}) bool {
	// check cache

	// check database
	mM.GORMPersister.FindInto(m, query, args)
	if len(m.GetUUID()) > 0 {
		return true
	}

	return false
}

// Delete - Deletes a model from the storages
func (mM ModelManager) Delete(m models.Model) bool {
	return mM.GORMPersister.Delete(m)
}
