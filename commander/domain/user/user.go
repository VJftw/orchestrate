package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User - A user of the application
type User struct {
	ID           uint      `gorm:"primary_key" json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	UUID         string    `json:"uuid" gorm:"unique"`
	EmailAddress string    `json:"emailAddress" gorm:"not null;unique;index" valid:"email,required"`
	Password     string    `json:"-" gorm:"-" valid:"length(6|255),required"`
	PasswordHash []byte    `json:"-" gorm:"not null"`
	FirstName    string    `json:"firstName" valid:"alpha"`
	LastName     string    `json:"lastName" valid:"alpha"`
}

// GetUUID - Return the UUID
func (u User) GetUUID() string {
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
