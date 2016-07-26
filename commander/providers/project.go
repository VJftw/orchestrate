package providers

import "github.com/vjftw/orchestrate/commander/models"

type IProject interface {
	New() *models.Project
}

type Project struct {
}

func NewProject() *Project {
	return &Project{}
}

func (pP Project) New() *models.Project {
	return &models.Project{}
}
