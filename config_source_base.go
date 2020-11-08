package servlet

import (
	"fmt"
	"sync"
)

// ConfigSourceBase defines a base code of a config source instance.
type ConfigSourceBase struct {
	mutex   sync.Locker
	partial ConfigPartial
}

// Close method used to be compliant with the container Closable interface.
func (*ConfigSourceBase) Close() {}

// Has will check if the requested path is present in the source
// configuration content.
func (s *ConfigSourceBase) Has(path string) bool {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.partial.Has(path)
}

// Get will retrieve the value stored in the requested path present in the
// configuration content.
// If the path does not exists, then the value nil will be returned.
// This method will mostly be used by the config object to obtain the full
// content of the source to aggregate all the data into his internal storing
// partial instance.
func (s *ConfigSourceBase) Get(path string) interface{} {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.partial.Get(path)
}
