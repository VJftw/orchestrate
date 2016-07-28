package project

import "github.com/vjftw/orchestrate/commander/persisters"

type Manager interface {
	Save(*Project) error
	GetInto(*Project, interface{}, ...interface{}) error
	Delete(*Project) error
}

type ProjectManager struct {
	GORMPersister persisters.Persister `inject:"persister.gorm"`
}

func NewManager() Manager {
	return &ProjectManager{}
}

// Save - Saves the model across storages
func (d ProjectManager) Save(u *Project) error {
	d.GORMPersister.Save(u)
	return nil
}

// GetInto - Searches the storages for a model identified by the query and places it into the given model reference.
// Returns true if found, false otherwise
func (d ProjectManager) GetInto(u *Project, query interface{}, args ...interface{}) error {
	// check cache

	// check database
	err := d.GORMPersister.GetInto(u, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Delete - Deletes a model from the storages
func (d ProjectManager) Delete(u *Project) error {
	d.GORMPersister.Delete(u)
	return nil
}
