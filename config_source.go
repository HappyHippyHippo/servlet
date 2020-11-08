package servlet

// ConfigSource defines the base interface of a config source.
type ConfigSource interface {
	Close()
	Has(path string) bool
	Get(path string) interface{}
}
