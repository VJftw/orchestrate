package user

import (
	"encoding/json"
	"errors"
	"io"
)

type Resolver interface {
	FromRequest(*User, io.ReadCloser) error
}

type UserResolver struct {
}

func (r UserResolver) FromRequest(u *User, b io.ReadCloser) error {
	var rJSON map[string]interface{}

	err := json.NewDecoder(b).Decode(&rJSON)
	if err != nil {
		return err
	}

	if _, ok := rJSON["emailAddress"]; !ok {
		return errors.New("Missing emailAddress")
	}

	if _, ok := rJSON["password"]; !ok {
		return errors.New("Missing password")
	}

	u.EmailAddress = rJSON["emailAddress"].(string)
	u.Password = rJSON["password"].(string)

	if val, ok := rJSON["firstName"]; ok {
		u.FirstName = val.(string)
	}

	if val, ok := rJSON["lastName"]; ok {
		u.LastName = val.(string)
	}

	return nil
}
