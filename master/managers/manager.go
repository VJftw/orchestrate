package managers

type IManager interface {
	Validate(v interface{}) error
	Save(v interface{})
}
