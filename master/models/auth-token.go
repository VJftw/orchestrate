package models

type AuthToken struct {
	Token string
}

func (aT AuthToken) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"authToken": aT.Token,
	}
}
