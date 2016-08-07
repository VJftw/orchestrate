package cadetGroup

import "time"

// CadetGroup - A group of Cadets which run the same container configuration
type CadetGroup struct {
	ID            uint      `gorm:"primary_key" json:"-"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	ProjectID     uint      `json:"-"`
	UUID          string    `json:"uuid" gorm:"unique"`
	Name          string    `json:"name" gorm:"not null" valid:"length(3,255),required"`
	Key           string    `json:"key" gorm:"not null"`
	Configuration string    `json:"configuration" valid:"required"`
}

// GetUUID - Returns the UUID
func (cG CadetGroup) GetUUID() string {
	return cG.UUID
}
