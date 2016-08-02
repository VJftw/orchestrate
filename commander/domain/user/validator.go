package user

import "github.com/asaskevich/govalidator"

type Validator interface {
	Validate(*User) bool
}

type UserValidator struct {
	UserManager Manager `inject:"user.manager"`
}

func (v UserValidator) Validate(u *User) bool {
	res, _ := govalidator.ValidateStruct(u)
	if !res {
		return false
	}

	_, err := v.UserManager.FindByEmailAddress(u.EmailAddress)
	if err == nil {
		return false
	}

	return true
}
