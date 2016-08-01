package cadetGroup

import (
	"github.com/jinzhu/gorm"
	"github.com/vjftw/orchestrate/commander/domain/project"
)

type Manager interface {
	NewForProject(*project.Project) *CadetGroup
	Save(*CadetGroup) error
	Delete(*CadetGroup) error
}

type CadetGroupManager struct {
	gorm *gorm.DB
}

func NewManager(gormDB *gorm.DB) Manager {
	return &CadetGroupManager{
		gorm: gormDB,
	}
}

func (m CadetGroupManager) NewForProject(project *project.Project) *CadetGroup {
	return &CadetGroup{
		ProjectID: project.ID,
	}
}

func (m CadetGroupManager) Save(c *CadetGroup) error {
	m.gorm.Save(c)
	return nil
}

func (m CadetGroupManager) Delete(c *CadetGroup) error {
	m.gorm.Delete(c)
	return nil
}
