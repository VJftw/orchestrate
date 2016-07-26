package messages

type ValidationMessage struct {
	Valid  bool        `json:"valid"`
	Errors interface{} `json:"errors"`
}

func (vM ValidationMessage) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"valid":  vM.Valid,
		"errors": vM.Errors,
	}
}
