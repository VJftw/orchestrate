package cadetGroup

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/vjftw/orchestrate/commander/domain/project"
	"github.com/vjftw/orchestrate/commander/domain/user"
)

// Manager - Manages CadetGroups
type Manager interface {
	NewForProject(*project.Project) *CadetGroup
	Save(*CadetGroup) error
	Delete(*CadetGroup) error
	FindByUserAndProjectUUID(*user.User, string) *[]CadetGroup
	FindByKey(string) (*CadetGroup, error)
}

type cadetGroupManager struct {
	gorm *gorm.DB
}

// NewManager - Returns a new Manager with the given GormDB
func NewManager(gormDB *gorm.DB) Manager {
	return &cadetGroupManager{
		gorm: gormDB,
	}
}

func (m cadetGroupManager) NewForProject(project *project.Project) *CadetGroup {
	return &CadetGroup{
		ProjectID: project.ID,
	}
}

func (m cadetGroupManager) Save(c *CadetGroup) error {
	m.gorm.Save(c)
	return nil
}

func (m cadetGroupManager) Delete(c *CadetGroup) error {
	m.gorm.Delete(c)
	return nil
}

func (m cadetGroupManager) FindByUserAndProjectUUID(user *user.User, projectUUID string) *[]CadetGroup {
	cadetGroups := []CadetGroup{}
	m.gorm.
		Table("cadet_groups").
		Joins("inner join projects on cadet_groups.project_id = projects.id").
		Joins("inner join users on projects.user_id = users.id").
		Where("projects.uuid = ? and user.id = ?", projectUUID, user.ID).
		Find(&cadetGroups)

	return &cadetGroups
}

func (m cadetGroupManager) FindByKey(key string) (*CadetGroup, error) {
	cG := &CadetGroup{}

	m.gorm.Where("key = ?", key).First(cG)

	if len(cG.GetUUID()) < 1 {
		return nil, errors.New("Not found")
	}

	return cG, nil
}
