package persisters

// Persister - interface defining what we can do with persistable structs
type Persister interface {
	Save(Persistable) error
	GetInto(Persistable, interface{}, ...interface{}) error
	Delete(Persistable) error
}

// Persistable - interface defining what is persistable
type Persistable interface {
	GetUUID() string
}
