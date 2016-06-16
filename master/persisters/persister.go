package persisters

// Persister - Persistence functions
type Persister interface {
	Save(interface{})
}
