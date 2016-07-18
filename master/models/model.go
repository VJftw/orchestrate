package models

// IModel -
type IModel interface {
	GetUUID() string
}

// ISerializable -
type ISerializable interface {
	ToMap() map[string]interface{}
}
