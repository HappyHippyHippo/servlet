package config

import (
	"fmt"
)

// DecoderFactory interface defines the methods of the decoder factory
// that can be used to instantiate a decoder.
type DecoderFactory interface {
	Register(strategy DecoderFactoryStrategy) error
	Create(format string, args ...interface{}) (Decoder, error)
}

type decoderFactory struct {
	strategies []DecoderFactoryStrategy
}

// NewDecoderFactory intantiate a new decoder factory.
func NewDecoderFactory() DecoderFactory {
	return &decoderFactory{
		strategies: []DecoderFactoryStrategy{},
	}
}

// Register will stores a new decoder factory strategy to be used
// to evalutate a request of a instance capable to parse a specific format.
// If the strategy accepts the format, then it will be used to instantiate the
// appropriate decoder that will be used to decode the configuration content.
func (f *decoderFactory) Register(strategy DecoderFactoryStrategy) error {
	if strategy == nil {
		return fmt.Errorf("Invalid nil 'strategy' argument")
	}

	f.strategies = append([]DecoderFactoryStrategy{strategy}, f.strategies...)

	return nil
}

// Create will instantiate the requested new decoder capable to
// parse the formatted content into a usable configuration partial.
func (f decoderFactory) Create(format string, args ...interface{}) (Decoder, error) {
	for _, s := range f.strategies {
		if s.Accept(format, args...) {
			return s.Create(args...)
		}
	}
	return nil, fmt.Errorf("Unrecognized format type : %s", format)
}
