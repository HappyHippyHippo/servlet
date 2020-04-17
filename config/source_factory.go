package config

import (
	"fmt"
)

// SourceFactory interface defines the methods of the source factory
// that can be used to instantiate a source of configuration content.
type SourceFactory interface {
	Register(strategy SourceFactoryStrategy) error
	Create(format string, args ...interface{}) (Source, error)
	CreateConfig(conf Partial) (Source, error)
}

type sourceFactory struct {
	strategies []SourceFactoryStrategy
}

// NewSourceFactory intantiate a new source factory.
func NewSourceFactory() SourceFactory {
	return &sourceFactory{
		strategies: []SourceFactoryStrategy{},
	}
}

// Register will register a new source factory strategy to be used
// on creation request.
func (f *sourceFactory) Register(strategy SourceFactoryStrategy) error {
	if strategy == nil {
		return fmt.Errorf("Invalid nil 'strategy' argument")
	}

	f.strategies = append([]SourceFactoryStrategy{strategy}, f.strategies...)

	return nil
}

// Create will instantiate and return a new config source by the type requested.
func (f sourceFactory) Create(stype string, args ...interface{}) (Source, error) {
	for _, s := range f.strategies {
		if s.Accept(stype, args...) {
			return s.Create(args...)
		}
	}
	return nil, fmt.Errorf("Unrecognized source type : %s", stype)
}

// CreateConfig will instantiate and return a new config source where the
// data used to decide the strategy to be used and also the initialization data
// comes frmo a configuration storing partial instance.
func (f sourceFactory) CreateConfig(conf Partial) (Source, error) {
	for _, s := range f.strategies {
		if s.AcceptConfig(conf) {
			return s.CreateConfig(conf)
		}
	}
	return nil, fmt.Errorf("Unrecognized source config : %v", conf)
}
