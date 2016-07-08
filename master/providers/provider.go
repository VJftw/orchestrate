package providers

// Provider - interface for Providers
type Provider interface {
	CreateNew() interface{}
}
