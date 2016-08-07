package cadetGroup

import "github.com/asaskevich/govalidator"

// Validator - Validates CadetGroups
type Validator interface {
	Validate(*CadetGroup) bool
}

type cadetGroupValidator struct {
}

// NewValidator - Returns a new Validator
func NewValidator() Validator {
	return &cadetGroupValidator{}
}

func (v cadetGroupValidator) Validate(c *CadetGroup) bool {
	res, _ := govalidator.ValidateStruct(c)

	return res
}
