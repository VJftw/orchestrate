package user

type Validator interface {
	Validate(*User) bool
}

type UserValidator struct {
}

func NewValidator() Validator {
	return &UserValidator{}
}

func (v UserValidator) Validate(u *User) bool {
	return true
}
