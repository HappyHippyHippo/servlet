package config

import (
	"io"
)

const (
	// DecoderFormatYAML defines the value to be used to declare a YAML
	// config source format.
	DecoderFormatYAML = "yaml"
)

type yamlDecoderFactoryStrategy struct{}

// NewYamlDecoderFactoryStrategy instantiate a new yaml decoder factory
// strategy that will enable the decoder factory to instantiate a new yaml
// decoder.
func NewYamlDecoderFactoryStrategy() DecoderFactoryStrategy {
	return &yamlDecoderFactoryStrategy{}
}

// Accept will check if the decoder factory strategy can instantiate a
// decoder giving the format and the creation request paramaters.
func (yamlDecoderFactoryStrategy) Accept(format string, args ...interface{}) bool {
	if format != DecoderFormatYAML || len(args) < 1 {
		return false
	}

	switch args[0].(type) {
	case io.Reader:
	default:
		return false
	}

	return true
}

// Create will instantiate the desired decoder instance with the given reader
// instance as source of the content to decode.
func (yamlDecoderFactoryStrategy) Create(args ...interface{}) (Decoder, error) {
	reader := args[0].(io.Reader)

	return NewYamlDecoder(reader)
}
