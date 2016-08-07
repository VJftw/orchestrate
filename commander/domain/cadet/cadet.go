package cadet

import "time"

// Cadet - a node to run containers on
type Cadet struct {
	ID           uint      `json:"-" gorm:"primary_key"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	CadetGroupID uint      `json:"-"`
	UUID         string    `json:"uuid" gorm:"unique"`
	Key          string    `json:"key" gorm:"not null"`
}

// GetUUID - Returns the UUID
func (c Cadet) GetUUID() string {
	return c.UUID
}
