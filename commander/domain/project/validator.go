package project

type Validator interface {
	Validate(*Project) bool
}

type ProjectValidator struct {
}

func NewValidator() Validator {
	return &ProjectValidator{}
}

func (v ProjectValidator) Validate(u *Project) bool {
	return true
}
