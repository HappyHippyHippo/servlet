package log

import (
	"fmt"
)

// FormatterFactory interface defines the methods of the formatter factory
// that can be used to instantiate a formatter.
type FormatterFactory interface {
	Register(strategy FormatterFactoryStrategy) error
	Create(format string, args ...interface{}) (Formatter, error)
}

type formatterFactory struct {
	strategies []FormatterFactoryStrategy
}

// NewFormatterFactory intantiate a new formatter factory.
func NewFormatterFactory() FormatterFactory {
	return &formatterFactory{
		strategies: []FormatterFactoryStrategy{},
	}
}

// Register will register a new formatter factory strategy to be used
// on creation request.
func (f *formatterFactory) Register(strategy FormatterFactoryStrategy) error {
	if strategy == nil {
		return fmt.Errorf("Invalid nil 'strategy' argument")
	}

	f.strategies = append([]FormatterFactoryStrategy{strategy}, f.strategies...)

	return nil
}

// Create will instantiate and return a new content formatter.
func (f formatterFactory) Create(format string, args ...interface{}) (Formatter, error) {
	for _, s := range f.strategies {
		if s.Accept(format, args...) {
			return s.Create(args...)
		}
	}
	return nil, fmt.Errorf("Unrecognized format type : %s", format)
}
