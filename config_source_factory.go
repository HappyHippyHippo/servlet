package servlet

import "fmt"

// ConfigSourceFactory defines a config source factory that uses a list of
// registered instantiation strategies to perform the config source
// instantiation.
type ConfigSourceFactory struct {
	strategies []ConfigSourceFactoryStrategy
}

// NewConfigSourceFactory instantiate a new source factory.
func NewConfigSourceFactory() *ConfigSourceFactory {
	return &ConfigSourceFactory{
		strategies: []ConfigSourceFactoryStrategy{},
	}
}

// Register will register a new source factory strategy to be used
// on creation request.
func (f *ConfigSourceFactory) Register(strategy ConfigSourceFactoryStrategy) error {
	if f == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if strategy == nil {
		return fmt.Errorf("invalid nil 'strategy' argument")
	}

	f.strategies = append([]ConfigSourceFactoryStrategy{strategy}, f.strategies...)

	return nil
}

// Create will instantiate and return a new config source by the type requested.
func (f ConfigSourceFactory) Create(sourceType string, args ...interface{}) (ConfigSource, error) {
	for _, s := range f.strategies {
		if s.Accept(sourceType, args...) {
			return s.Create(args...)
		}
	}
	return nil, fmt.Errorf("unrecognized source type : %s", sourceType)
}

// CreateConfig will instantiate and return a new config source where the
// data used to decide the strategy to be used and also the initialization
// data comes from a configuration storing partial instance.
func (f ConfigSourceFactory) CreateConfig(conf ConfigPartial) (ConfigSource, error) {
	for _, s := range f.strategies {
		if s.AcceptConfig(conf) {
			return s.CreateConfig(conf)
		}
	}
	return nil, fmt.Errorf("unrecognized source config : %v", conf)
}
