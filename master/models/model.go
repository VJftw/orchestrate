package models

// IModel -
type IModel interface {
}

// ISerializable -
type ISerializable interface {
	ToMap() map[string]interface{}
}
