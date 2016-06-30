package models

// IModel -
type IModel interface {
}

// Serializable -
type Serializable interface {
	ToMap() map[string]interface{}
}
