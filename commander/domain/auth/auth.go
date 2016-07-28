package auth

type Auth struct {
	Token string
}

func (a Auth) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"authToken": a.Token,
	}
}
