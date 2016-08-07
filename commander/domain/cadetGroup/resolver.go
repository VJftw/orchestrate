package cadetGroup

import (
	"encoding/json"
	"errors"
	"io"
)

// Resolver - Resolves data into CadetGroups
type Resolver interface {
	FromRequest(*CadetGroup, io.ReadCloser) error
}

type cadetGroupResolver struct {
}

// NewResolver - returns a new Resolver
func NewResolver() Resolver {
	return &cadetGroupResolver{}
}

func (r cadetGroupResolver) FromRequest(c *CadetGroup, b io.ReadCloser) error {
	var rJSON map[string]interface{}

	err := json.NewDecoder(b).Decode(&rJSON)
	if err != nil {
		return err
	}

	if _, ok := rJSON["name"]; !ok {
		return errors.New("Missing name")
	}

	if _, ok := rJSON["configuration"]; !ok {
		return errors.New("Missing configuration")
	}

	c.Name = rJSON["name"].(string)
	c.Configuration = rJSON["configuration"].(string)

	return nil
}
