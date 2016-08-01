package cadetGroup

import (
	"encoding/json"
	"io"
)

type Resolver interface {
	FromRequest(*CadetGroup, io.ReadCloser) error
}

type CadetGroupResolver struct {
}

func NewResolver() Resolver {
	return &CadetGroupResolver{}
}

func (r CadetGroupResolver) FromRequest(c *CadetGroup, b io.ReadCloser) error {
	var rJSON map[string]interface{}

	err := json.NewDecoder(b).Decode(&rJSON)
	if err != nil {
		return err
	}

	c.Name = rJSON["name"].(string)
	c.Configuration = rJSON["configuration"].(string)

	return nil
}
