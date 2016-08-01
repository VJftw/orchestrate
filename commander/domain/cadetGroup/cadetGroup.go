package cadetGroup

import "time"

type CadetGroup struct {
	ID            uint      `gorm:"primary_key" json:"-"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	ProjectID     uint      `json:"-"`
	UUID          string    `json:"uuid" gorm:"unique"`
	Name          string    `json:"name" gorm:"not null" valid:"alpha, required"`
	Key           string    `json:"-" gorm:"not null"`
	Configuration string    `json:"configuration" valid:"alpha"`
}

// GetUUID - Returns the UUID
func (cG CadetGroup) GetUUID() string {
	return cG.UUID
}
