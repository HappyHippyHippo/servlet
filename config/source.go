package config

import "github.com/happyhippyhippo/servlet/sys"

// Source interface defines the methods to interact with a config source.
type Source interface {
	Close() error
	Has(path string) bool
	Get(path string) interface{}
}

// ObservableSource interface extends the Source interface with methods
// specific to sources that will be checked for updates.
type ObservableSource interface {
	Source
	Reload() (bool, error)
}

type source struct {
	mutex   sys.RWMutex
	partial Partial
}

// Close method used to be compliant with the container Closable interface.
func (*source) Close() error {
	return nil
}

// Has will check if the requested path is present in the source
// configuration path.
func (s *source) Has(path string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.partial.Has(path)
}

// Get will retrieve the value stored in the requested path present in the
// configuration source.
// If the path does not exists, then the value nil will be returned.
func (s *source) Get(path string) interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.partial.Get(path)
}
