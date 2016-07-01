package models

// IModel -
type IModel interface {
	GetID() uint
}

// Serializable -
type Serializable interface {
	ToMap() map[string]interface{}
}
