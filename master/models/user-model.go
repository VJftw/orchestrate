package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User - A user of the application
type User struct {
	gorm.Model
	UUID         []byte `json:"uuid" gorm:"unique"`
	EmailAddress string `json:"emailAddress" gorm:"not null;unique;index" valid:"email,required"`
	Password     string `json:"password" gorm:"-" valid:"length(6|255),required"`
	PasswordHash []byte `gorm:"non null"`
	FirstName    string `json:"firstName" valid:"alpha"`
	LastName     string `json:"lastName" valid:"alpha"`
}

// GetUUID - Return the UUID
func (u User) GetUUID() []byte {
	return u.UUID
}

// EncryptPassword - Sets the PasswordHash
func (u *User) EncryptPassword() {
	u.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(u.Password), 10)
}

// VerifyPassword - Verifies Password
func (u User) VerifyPassword() bool {
	if bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(u.Password)) == nil {
		return true
	}
	return false
}

// ToMap - Returns a map representation of a User
func (u User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"uuid":         string(u.UUID),
		"emailAddress": u.EmailAddress,
		"firstName":    u.FirstName,
		"lastName":     u.LastName,
	}
}
