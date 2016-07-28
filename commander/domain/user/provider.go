package user

type Provider interface {
	New() *User
}

type UserProvider struct {
}

func NewProvider() Provider {
	return &UserProvider{}
}

func (p UserProvider) New() *User {
	return &User{}
}
