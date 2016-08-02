package project

import (
	"encoding/json"
	"errors"
	"io"
)

type Resolver interface {
	FromRequest(*Project, io.ReadCloser) error
}

type ProjectResolver struct {
}

func (r ProjectResolver) FromRequest(p *Project, b io.ReadCloser) error {
	var rJSON map[string]interface{}

	err := json.NewDecoder(b).Decode(&rJSON)
	if err != nil {
		return err
	}

	if _, ok := rJSON["name"]; !ok {
		return errors.New("Missing name")
	}
	p.Name = rJSON["name"].(string)

	return nil
}
