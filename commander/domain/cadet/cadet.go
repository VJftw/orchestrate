package cadet

import "time"

// Cadet - a node to run containers on
type Cadet struct {
	ID           uint      `gorm:"primary_key" json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	CadetGroupID int       `json:"-"`
	UUID         string    `json:"uuid" gorm:"unique"`
}

// GetUUID - Returns the UUID
func (c Cadet) GetUUID() string {
	return c.UUID
}
