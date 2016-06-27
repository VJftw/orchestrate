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
