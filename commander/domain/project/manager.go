package project

import (
	"errors"

	"github.com/vjftw/orchestrate/commander/domain/user"

	"github.com/jinzhu/gorm"
)

type Manager interface {
	NewForUser(*user.User) *Project
	Save(*Project) error
	Delete(*Project) error
	GetInto(*Project, interface{}, ...interface{})
	FindByUserAndUUID(*user.User, string) (*Project, error)
	FindByUser(*user.User) *[]Project
}

type ProjectManager struct {
	gorm *gorm.DB
}

func NewManager(gormDB *gorm.DB) Manager {
	return &ProjectManager{
		gorm: gormDB,
	}
}

func (m ProjectManager) NewForUser(u *user.User) *Project {
	return &Project{
		UserID: u.ID,
	}
}

// Save - Saves the model across storages
func (m ProjectManager) Save(u *Project) error {
	m.gorm.Save(u)
	return nil
}

// GetInto - Searches the storages for a model identified by the query and places it into the given model reference.
// Returns true if found, false otherwise
func (m ProjectManager) GetInto(p *Project, query interface{}, args ...interface{}) {
	// check database
	m.gorm.Where(query, args...).First(p)
}

func (m ProjectManager) FindByUserAndUUID(u *user.User, UUID string) (*Project, error) {
	project := &Project{}
	m.gorm.
		Table("projects").
		Joins("inner join users on projects.user_id = users.id").
		Where("projects.uuid = ? and users.id = ?", UUID, u.ID).
		First(project)

	if len(project.GetUUID()) < 1 {
		return nil, errors.New("Not found")
	}
	return project, nil
}

func (m ProjectManager) FindByUser(u *user.User) *[]Project {
	projects := []Project{}
	m.gorm.Where("user_id = ?", u.ID).Find(&projects)

	return &projects
}

// Delete - Deletes a model from the storages
func (m ProjectManager) Delete(u *Project) error {
	m.gorm.Delete(u)
	return nil
}
