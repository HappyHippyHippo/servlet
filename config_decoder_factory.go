package servlet

import "fmt"

// ConfigDecoderFactory defined the instance used to instantiate a new config
// stream decoder for a specific encoding format.
type ConfigDecoderFactory struct {
	strategies []ConfigDecoderFactoryStrategy
}

// NewConfigDecoderFactory instantiate a new decoder factory.
func NewConfigDecoderFactory() *ConfigDecoderFactory {
	return &ConfigDecoderFactory{
		strategies: []ConfigDecoderFactoryStrategy{},
	}
}

// Register will stores a new decoder factory strategy to be used
// to evaluate a request of a instance capable to parse a specific format.
// If the strategy accepts the format, then it will be used to instantiate the
// appropriate decoder that will be used to decode the configuration content.
func (f *ConfigDecoderFactory) Register(strategy ConfigDecoderFactoryStrategy) error {
	if f == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if strategy == nil {
		return fmt.Errorf("invalid nil 'strategy' argument")
	}

	f.strategies = append([]ConfigDecoderFactoryStrategy{strategy}, f.strategies...)

	return nil
}

// Create will instantiate the requested new decoder capable to
// parse the formatted content into a usable configuration partial.
func (f ConfigDecoderFactory) Create(format string, args ...interface{}) (ConfigDecoder, error) {
	for _, s := range f.strategies {
		if s.Accept(format, args...) {
			return s.Create(args...)
		}
	}
	return nil, fmt.Errorf("unrecognized format type : %s", format)
}
