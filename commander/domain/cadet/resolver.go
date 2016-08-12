package cadet

import (
	"encoding/json"
	"errors"
	"io"
)

type Resolver interface {
	KeyFromRequest(io.ReadCloser) (string, error)
}

type CadetResolver struct {
}

func NewResolver() Resolver {
	return &CadetResolver{}
}

func (r CadetResolver) KeyFromRequest(b io.ReadCloser) (string, error) {
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
