package config

import (
	"fmt"

	"github.com/spf13/afero"
)

const (
	// SourceTypeObservableFile defines the value to be used to declare a
	// observable file config source type.
	SourceTypeObservableFile = "observable_file"
)

type observableFileSourceFactoryStrategy struct {
	fileSystem     afero.Fs
	decoderFactory DecoderFactory
}

// NewObservableFileSourceFactoryStrategy instantiate a new yaml decoder factory
// strategy that will enable the decoder factory to instantiate a new yaml
// decoder.
func NewObservableFileSourceFactoryStrategy(fileSystem afero.Fs, decoderFactory DecoderFactory) (SourceFactoryStrategy, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("Invalid nil 'fileSystem' argument")
	}
	if decoderFactory == nil {
		return nil, fmt.Errorf("Invalid nil 'decoderFactory' argument")
	}

	return &observableFileSourceFactoryStrategy{
		fileSystem:     fileSystem,
		decoderFactory: decoderFactory,
	}, nil
}

// Accept will check if the source factory strategy can instantiate a
// decoder of the requested format.
func (observableFileSourceFactoryStrategy) Accept(stype string, args ...interface{}) bool {
	if stype != SourceTypeObservableFile || len(args) < 2 {
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
// decoder where the data to check comes from a configuration instance.
func (s observableFileSourceFactoryStrategy) AcceptConfig(conf Partial) (check bool) {
	defer func() {
		if r := recover(); r != nil {
			check = false
		}
	}()

	stype := conf.String("type")
	path := conf.String("path")
	format := conf.String("format")

	return s.Accept(stype, path, format)
}

// Create will instantiate the desired source instance.
func (s observableFileSourceFactoryStrategy) Create(args ...interface{}) (Source, error) {
	path := args[0].(string)
	format := args[1].(string)

	return NewObservableFileSource(path, format, s.fileSystem, s.decoderFactory)
}

// CreateConfig will instantiate the desired source instance where the initialization
// data comes from a configuration instance.
func (s observableFileSourceFactoryStrategy) CreateConfig(conf Partial) (Source, error) {
	path := conf.String("path")
	format := conf.String("format")

	return s.Create(path, format)
}
