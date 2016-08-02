package user

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Manager interface {
	New() *User
	Save(*User) error
	GetInto(*User, interface{}, ...interface{})
	Delete(*User) error
	FindByUUID(string) (*User, error)
	FindByEmailAddress(string) (*User, error)
}

type UserManager struct {
	gorm *gorm.DB
}

func NewManager(gormDB *gorm.DB) Manager {
	return &UserManager{
		gorm: gormDB,
	}
}

func (m UserManager) New() *User {
	return &User{}
}

// Save - Saves the model across storages
func (m UserManager) Save(u *User) error {
	m.gorm.Save(u)
	return nil
}

// GetInto - Searches the storages for a model identified by the query and places it into the given model reference.
// Returns true if found, false otherwise
func (m UserManager) GetInto(u *User, query interface{}, args ...interface{}) {
	// check database
	m.gorm.Where(query, args...).First(u)
}

func (m UserManager) FindByUUID(uuid string) (*User, error) {
	user := User{}

	m.GetInto(&user, "uuid = ?", uuid)

	if len(user.GetUUID()) < 1 {
		return nil, errors.New("Not found")
	}

	return &user, nil
}

func (m UserManager) FindByEmailAddress(emailAddress string) (*User, error) {
	user := User{}

	m.GetInto(&user, "email_address = ?", emailAddress)

	if len(user.GetUUID()) < 1 {
		return nil, errors.New("Not found")
	}

	return &user, nil
}

// Delete - Deletes a model from the storages
func (m UserManager) Delete(u *User) error {
	m.gorm.Delete(u)
	return nil
}
