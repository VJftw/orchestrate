package project

import (
	"encoding/json"
	"io"
)

type Resolver interface {
	FromRequest(*Project, io.ReadCloser) error
}

type ProjectResolver struct {
}

func NewResolver() Resolver {
	return &ProjectResolver{}
}

func (r ProjectResolver) FromRequest(p *Project, b io.ReadCloser) error {
	var rJSON map[string]interface{}

	err := json.NewDecoder(b).Decode(&rJSON)
	if err != nil {
		return err
	}

	p.Name = rJSON["name"].(string)

	return nil
}
