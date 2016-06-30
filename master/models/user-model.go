package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User - A user of the application
type User struct {
	gorm.Model
	EmailAddress string `json:"emailAddress" gorm:"not null;unique" valid:"email,required"`
	Password     string `json:"password" gorm:"-" valid:"length(6|255),required"`
	PasswordHash []byte `gorm:"non null"`
	FirstName    string `json:"firstName" valid:"alpha"`
	LastName     string `json:"lastName" valid:"alpha"`
}

// EncryptPassord - Sets the PasswordHash
func (u *User) EncryptPassord() {
	u.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(u.Password), 10)
}

// ToMap - Returns a map representation of a User
func (u User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"emailAddress": u.EmailAddress,
		"firstName":    u.FirstName,
		"lastName":     u.LastName,
	}
}
