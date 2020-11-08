package servlet

import (
	"fmt"
	"github.com/spf13/afero"
)

// ConfigSourceFactoryStrategyObservableFile defines a observable config file
// source instantiation strategy to be used by the config sources factory
// instance.
type ConfigSourceFactoryStrategyObservableFile struct {
	fileSystem     afero.Fs
	decoderFactory *ConfigDecoderFactory
}

// NewConfigSourceFactoryStrategyObservableFile instantiate a new observable
// file source factory strategy that will enable the source factory to
// instantiate a new observable file configuration source.
func NewConfigSourceFactoryStrategyObservableFile(fileSystem afero.Fs, decoderFactory *ConfigDecoderFactory) (*ConfigSourceFactoryStrategyObservableFile, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("invalid nil 'fileSystem' argument")
	}
	if decoderFactory == nil {
		return nil, fmt.Errorf("invalid nil 'decoderFactory' argument")
	}

	return &ConfigSourceFactoryStrategyObservableFile{
		fileSystem:     fileSystem,
		decoderFactory: decoderFactory,
	}, nil
}

// Accept will check if the source factory strategy can instantiate a
// new source of the requested type. Also, validates that there is the path
// and content format extra parameters, and thar this parameters are strings.
func (ConfigSourceFactoryStrategyObservableFile) Accept(sourceType string, args ...interface{}) bool {
	if sourceType != ConfigSourceTypeObservableFile || len(args) < 2 {
		return false
	}

	switch args[0].(type) {
	case string:
	default:
		return false
	}

	switch args[1].(type) {
	case string:
	default:
		return false
	}

	return true
}

// AcceptConfig will check if the source factory strategy can instantiate a
// source where the data to check comes from a configuration partial instance.
func (s ConfigSourceFactoryStrategyObservableFile) AcceptConfig(conf ConfigPartial) (check bool) {
	defer func() {
		if r := recover(); r != nil {
			check = false
		}
	}()

	sourceType := conf.String("type")
	path := conf.String("path")
	format := conf.String("format")

	return s.Accept(sourceType, path, format)
}

// Create will instantiate the desired observable file source instance.
func (s ConfigSourceFactoryStrategyObservableFile) Create(args ...interface{}) (source ConfigSource, err error) {
	defer func() {
		if r := recover(); r != nil {
			source = nil
			err = r.(error)
		}
	}()

	path := args[0].(string)
	format := args[1].(string)

	return NewConfigSourceObservableFile(path, format, s.fileSystem, s.decoderFactory)
}

// CreateConfig will instantiate the desired observable file source instance
// where the initialization data comes from a configuration partial instance.
func (s ConfigSourceFactoryStrategyObservableFile) CreateConfig(conf ConfigPartial) (source ConfigSource, err error) {
	defer func() {
		if r := recover(); r != nil {
			source = nil
			err = r.(error)
		}
	}()

	path := conf.String("path")
	format := conf.String("format")

	return s.Create(path, format)
}
