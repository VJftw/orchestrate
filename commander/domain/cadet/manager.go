package cadet

import (
	"github.com/jinzhu/gorm"
	"github.com/vjftw/orchestrate/commander/domain/cadetGroup"
)

type Manager interface {
	NewForCadetGroup(*cadetGroup.CadetGroup) *Cadet
	Save(*Cadet) error
	Delete(*Cadet) error
}

type CadetManager struct {
	gorm *gorm.DB
}

func NewManager(gormDB *gorm.DB) Manager {
	return &CadetManager{
		gorm: gormDB,
	}
}

func (m CadetManager) NewForCadetGroup(cadetGroup *cadetGroup.CadetGroup) *Cadet {
	return &Cadet{
		CadetGroupID: cadetGroup.ID,
	}
}

func (m CadetManager) Save(c *Cadet) error {
	m.gorm.Save(c)
	return nil
}

func (m CadetManager) Delete(c *Cadet) error {
	m.gorm.Delete(c)
	return nil
}
