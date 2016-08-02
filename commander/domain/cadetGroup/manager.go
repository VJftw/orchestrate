package cadetGroup

import (
	"github.com/jinzhu/gorm"
	"github.com/vjftw/orchestrate/commander/domain/project"
)

type Manager interface {
	NewForProject(*project.Project) *CadetGroup
	Save(*CadetGroup) error
	Delete(*CadetGroup) error
	FindByProject(*project.Project) *[]CadetGroup
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

func (m CadetGroupManager) FindByProject(p *project.Project) *[]CadetGroup {
	cadetGroups := []CadetGroup{}
	m.gorm.
		Table("cadet_groups").Joins("inner join projects on cadet_groups.project_id = projects.id").
		Where("projects.id = ?", p.ID).
		Find(&cadetGroups)

	return &cadetGroups
}
