package cadetGroup

type Validator interface {
	Validate(*CadetGroup) bool
}

type CadetGroupValidator struct {
}

func NewValidator() Validator {
	return &CadetGroupValidator{}
}

func (v CadetGroupValidator) Validate(c *CadetGroup) bool {
	return true
}
