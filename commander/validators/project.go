package validators

import "github.com/vjftw/orchestrate/commander/models"

type IProject interface {
	Validate(*models.Project) bool
}

type Project struct {
}

func NewProject() *Project {
	return &Project{}
}

// Validate - Validates a Project model
func (pV Project) Validate(p *models.Project) bool {
	return true
}
