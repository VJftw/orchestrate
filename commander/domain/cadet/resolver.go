package cadet

import (
	"encoding/json"
	"errors"
	"io"
)

// Resolver - interface defining methods for a Cadet Resolver
type Resolver interface {
	KeyFromRequest(io.ReadCloser) (string, error)
}

type cadetResolver struct {
}

// NewResolver - Returns a new Cadet Resolver
func NewResolver() Resolver {
	return &cadetResolver{}
}

func (r cadetResolver) KeyFromRequest(b io.ReadCloser) (string, error) {
	var rJSON map[string]interface{}

	err := json.NewDecoder(b).Decode(&rJSON)
	if err != nil {
		return "", err
	}

	if _, ok := rJSON["cadetGroupKey"]; !ok {
		return "", errors.New("Missing cadet group key")
	}
	cadetGroupKey := rJSON["cadetGroupKey"].(string)

	return cadetGroupKey, nil
}
