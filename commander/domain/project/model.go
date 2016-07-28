package project

import "github.com/jinzhu/gorm"

// Project - A Project belonging to a User
type Project struct {
	gorm.Model
	UserID int
	UUID   string `json:"uuid" gorm:"unique"`
	Name   string `json:"name" gorm:"not null" valid:"alpha,required"`
}

// GetUUID - Returns the UUID
func (p Project) GetUUID() string {
	return p.UUID
}

// ToMap - Returns a map representation of a Project
func (p Project) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"uuid": p.UUID,
		"name": p.Name,
	}
}
