package models

import "github.com/jinzhu/gorm"

// User - A user of the application
type User struct {
	gorm.Model
	EmailAddress string `json:"emailAddress" gorm:"not null;unique" valid:"email,required"`
	Password     string `json:"password" gorm:"not null;" valid:"length(6|255),required"`
	FirstName    string `json:"firstName" valid:"alpha"`
	LastName     string `json:"lastName" valid:"alpha"`
}

// ToMap - Returns a map representation of a User
func (u User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"emailAddress": u.EmailAddress,
		"firstName":    u.FirstName,
		"lastName":     u.LastName,
	}
}
