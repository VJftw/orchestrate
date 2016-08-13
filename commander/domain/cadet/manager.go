package cadet

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/vjftw/orchestrate/commander/domain/cadetGroup"
)

// Manager - interface defining methods for a Cadet Manager
type Manager interface {
	NewForCadetGroup(*cadetGroup.CadetGroup) *Cadet
	Save(*Cadet) error
	Delete(*Cadet) error
	FindByUUID(string) (*Cadet, error)
}

type cadetManager struct {
	gorm *gorm.DB
}

// NewManager - Returns a new Manager
func NewManager(gormDB *gorm.DB) Manager {
	return &cadetManager{
		gorm: gormDB,
	}
}

func (m cadetManager) NewForCadetGroup(cadetGroup *cadetGroup.CadetGroup) *Cadet {
	return &Cadet{
		CadetGroupID: cadetGroup.ID,
	}
}

func (m cadetManager) Save(c *Cadet) error {
	m.gorm.Save(c)
	return nil
}

func (m cadetManager) Delete(c *Cadet) error {
	m.gorm.Delete(c)
	return nil
}

func (m cadetManager) FindByUUID(cadetUUID string) (*Cadet, error) {
	cadet := &Cadet{}

	m.gorm.Where("uuid = ?", cadetUUID).First(cadet)

	if len(cadet.GetUUID()) < 1 {
		return nil, errors.New("Not found")
	}

	return cadet, nil
}
