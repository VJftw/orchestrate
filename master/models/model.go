package models

// Model -
type Model interface {
	GetUUID() []byte
}

// Serializable -
type Serializable interface {
	ToMap() map[string]interface{}
}
