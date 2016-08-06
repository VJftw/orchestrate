package project

import "github.com/asaskevich/govalidator"

type Validator interface {
	Validate(*Project) bool
}

type ProjectValidator struct {
}

func NewValidator() Validator {
	return &ProjectValidator{}
}

func (v ProjectValidator) Validate(p *Project) bool {
	res, _ := govalidator.ValidateStruct(p)

	return res
}
