package persisters

// IPersister - Persistence functions
type IPersister interface {
	Save(interface{})
	FindInto(interface{}, interface{}, ...interface{})
}
