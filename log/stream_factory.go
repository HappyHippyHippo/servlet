package log

import (
	"fmt"

	"github.com/happyhippyhippo/servlet/config"
)

// StreamFactory interface defines the methods of the stream factory
// that can be used to instantiate a stream.
type StreamFactory interface {
	Register(strategy StreamFactoryStrategy) error
	Create(format string, args ...interface{}) (Stream, error)
	CreateConfig(conf config.Partial) (Stream, error)
}

type streamFactory struct {
	strategies []StreamFactoryStrategy
}

// NewStreamFactory intantiate a new stream factory.
func NewStreamFactory() StreamFactory {
	return &streamFactory{
		strategies: []StreamFactoryStrategy{},
	}
}

// Register will register a new stream factory strategy to be used
// on creation request.
func (f *streamFactory) Register(strategy StreamFactoryStrategy) error {
	if strategy == nil {
		return fmt.Errorf("Invalid nil 'strategy' argument")
	}

	f.strategies = append([]StreamFactoryStrategy{strategy}, f.strategies...)

	return nil
}

// Create will instantiate and return a new config stream.
func (f streamFactory) Create(stype string, args ...interface{}) (Stream, error) {
	for _, s := range f.strategies {
		if s.Accept(stype, args...) {
			return s.Create(args...)
		}
	}
	return nil, fmt.Errorf("Unrecognized stream type : %s", stype)
}

// CreateConfig will instantiate and return a new config stream loaded by a
// configuration instance.
func (f streamFactory) CreateConfig(conf config.Partial) (Stream, error) {
	for _, s := range f.strategies {
		if s.AcceptConfig(conf) {
			return s.CreateConfig(conf)
		}
	}
	return nil, fmt.Errorf("Unrecognized stream config : %v", conf)
}
