package models

import "github.com/jinzhu/gorm"

// User - A user of the application
type User struct {
	gorm.Model
	ID           int    `json:"id" gorm:"AUTO_INCREMENT"`
	EmailAddress string `json:"emailAddress" gorm:"not null;unique"`
	Password     string `json:"password" gorm:"not null;"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}
