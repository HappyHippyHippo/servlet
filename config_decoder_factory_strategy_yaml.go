package servlet

import "io"

// ConfigDecoderFactoryStrategyYaml defines a strategy used to instantiate
// a YAML config stream decoder.
type ConfigDecoderFactoryStrategyYaml struct{}

// NewConfigDecoderFactoryStrategyYaml instantiate a new yaml decoder factory
// strategy that will enable the decoder factory to instantiate a new yaml
// decoder.
func NewConfigDecoderFactoryStrategyYaml() *ConfigDecoderFactoryStrategyYaml {
	return &ConfigDecoderFactoryStrategyYaml{}
}

// Accept will check if the decoder factory strategy can instantiate a
// decoder giving the format and the creation request parameters.
func (ConfigDecoderFactoryStrategyYaml) Accept(format string, args ...interface{}) bool {
	if format != ConfigDecoderFormatYAML || len(args) < 1 {
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
func (ConfigDecoderFactoryStrategyYaml) Create(args ...interface{}) (ConfigDecoder, error) {
	reader := args[0].(io.Reader)

	return NewConfigDecoderYaml(reader)
}
