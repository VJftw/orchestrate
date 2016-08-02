package user

import "github.com/asaskevich/govalidator"

type Validator interface {
	Validate(*User) bool
}

type UserValidator struct {
}

func NewValidator() Validator {
	return &UserValidator{}
}

func (v UserValidator) Validate(u *User) bool {
	res, _ := govalidator.ValidateStruct(u)

	return res
}
