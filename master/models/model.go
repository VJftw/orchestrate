package models

// Model -
type Model interface {
	GetUUID() string
}

// Serializable -
type Serializable interface {
	ToMap() map[string]interface{}
}
