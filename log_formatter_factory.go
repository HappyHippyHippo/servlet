package servlet

import "fmt"

// LogFormatterFactory defines the log formatter factory structure used to
// instantiate log formatters, based on registered instantiation strategies.
type LogFormatterFactory struct {
	strategies []LogFormatterFactoryStrategy
}

// NewLogFormatterFactory instantiate a new formatter factory.
func NewLogFormatterFactory() *LogFormatterFactory {
	return &LogFormatterFactory{
		strategies: []LogFormatterFactoryStrategy{},
	}
}

// Register will register a new formatter factory strategy to be used
// on requesting to create a formatter for a defined format.
func (f *LogFormatterFactory) Register(strategy LogFormatterFactoryStrategy) error {
	if f == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if strategy == nil {
		return fmt.Errorf("invalid nil 'strategy' argument")
	}

	f.strategies = append([]LogFormatterFactoryStrategy{strategy}, f.strategies...)

	return nil
}

// Create will instantiate and return a new content formatter.
func (f LogFormatterFactory) Create(format string, args ...interface{}) (LogFormatter, error) {
	for _, s := range f.strategies {
		if s.Accept(format, args...) {
			return s.Create(args...)
		}
	}
	return nil, fmt.Errorf("unrecognized format type : %s", format)
}
