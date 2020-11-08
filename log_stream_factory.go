package servlet

import "fmt"

// LogStreamFactory defines a log stream factory instance used to instantiate
// log stream based in the registered log stream instantiation strategies.
type LogStreamFactory struct {
	strategies []LogStreamFactoryStrategy
}

// NewLogStreamFactory instantiate a new stream factory.
func NewLogStreamFactory() *LogStreamFactory {
	return &LogStreamFactory{
		strategies: []LogStreamFactoryStrategy{},
	}
}

// Register will register a new stream factory strategy to be used
// on creation requests.
func (f *LogStreamFactory) Register(strategy LogStreamFactoryStrategy) error {
	if f == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if strategy == nil {
		return fmt.Errorf("invalid nil 'strategy' argument")
	}

	f.strategies = append([]LogStreamFactoryStrategy{strategy}, f.strategies...)

	return nil
}

// Create will instantiate and return a new config stream.
func (f LogStreamFactory) Create(sourceType string, args ...interface{}) (LogStream, error) {
	for _, s := range f.strategies {
		if s.Accept(sourceType, args...) {
			return s.Create(args...)
		}
	}
	return nil, fmt.Errorf("unrecognized stream type : %s", sourceType)
}

// CreateConfig will instantiate and return a new config stream loaded by a
// configuration instance.
func (f LogStreamFactory) CreateConfig(conf ConfigPartial) (LogStream, error) {
	for _, s := range f.strategies {
		if s.AcceptConfig(conf) {
			return s.CreateConfig(conf)
		}
	}
	return nil, fmt.Errorf("unrecognized stream config : %v", conf)
}
