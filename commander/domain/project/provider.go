package project

type Provider interface {
	New() *Project
}

type ProjectProvider struct {
}

func NewProvider() Provider {
	return &ProjectProvider{}
}

func (p ProjectProvider) New() *Project {
	return &Project{}
}
