package user

import "github.com/vjftw/orchestrate/commander/persisters"

type Manager interface {
	Save(*User) error
	GetInto(*User, interface{}, ...interface{}) error
	Delete(*User) error
}

type UserManager struct {
	GORMPersister persisters.Persister `inject:"persister.gorm"`
}

func NewManager() Manager {
	return &UserManager{}
}

// Save - Saves the model across storages
func (d UserManager) Save(u *User) error {
	d.GORMPersister.Save(u)
	return nil
}

// GetInto - Searches the storages for a model identified by the query and places it into the given model reference.
// Returns true if found, false otherwise
func (d UserManager) GetInto(u *User, query interface{}, args ...interface{}) error {
	// check cache

	// check database
	err := d.GORMPersister.GetInto(u, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Delete - Deletes a model from the storages
func (d UserManager) Delete(u *User) error {
	d.GORMPersister.Delete(u)
	return nil
}
