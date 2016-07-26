package resolvers

import (
	"encoding/json"
	"io"

	"github.com/vjftw/orchestrate/commander/models"
)

type IProject interface {
	FromRequest(*models.Project, io.ReadCloser) error
}

type Project struct {
}

func NewProject() *Project {
	return &Project{}
}

func (pR Project) FromRequest(p *models.Project, b io.ReadCloser) error {
	var rJSON map[string]interface{}

	err := json.NewDecoder(b).Decode(&rJSON)
	if err != nil {
		return err
	}

	p.Name = rJSON["name"].(string)

	return nil
}
