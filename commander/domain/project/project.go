package project

import "time"

// Project - A Project belonging to a User
type Project struct {
	ID        uint      `gorm:"primary_key" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	UserID    uint
	UUID      string `json:"uuid" gorm:"unique"`
	Name      string `json:"name" gorm:"not null" valid:"alpha,required"`
}

// GetUUID - Returns the UUID
func (p Project) GetUUID() string {
	return p.UUID
}
