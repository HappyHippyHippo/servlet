package config

import (
	"fmt"

	"github.com/spf13/afero"
)

const (
	// SourceTypeFile defines the value to be used to declare a simple file
	// config source type.
	SourceTypeFile = "file"
)

type fileSourceFactoryStrategy struct {
	fileSystem     afero.Fs
	decoderFactory DecoderFactory
}

// NewFileSourceFactoryStrategy instantiate a new file source factory
// strategy that will enable the source factory to instantiate a new
// file configuration source.
func NewFileSourceFactoryStrategy(fileSystem afero.Fs, decoderFactory DecoderFactory) (SourceFactoryStrategy, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("Invalid nil 'fileSystem' argument")
	}
	if decoderFactory == nil {
		return nil, fmt.Errorf("Invalid nil 'decoderFactory' argument")
	}

	return &fileSourceFactoryStrategy{
		fileSystem:     fileSystem,
		decoderFactory: decoderFactory,
	}, nil
}

// Accept will check if the source factory strategy can instantiate a
// new source of the requested type. Also, validates that there is the path
// and content format extra parameters, and thar this parameters are strings.
func (fileSourceFactoryStrategy) Accept(stype string, args ...interface{}) bool {
	if stype != SourceTypeFile || len(args) < 2 {
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
func (s fileSourceFactoryStrategy) AcceptConfig(conf Partial) (check bool) {
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

// Create will instantiate the desired file source instance.
func (s fileSourceFactoryStrategy) Create(args ...interface{}) (Source, error) {
	path := args[0].(string)
	format := args[1].(string)

	return NewFileSource(path, format, s.fileSystem, s.decoderFactory)
}

// CreateConfig will instantiate the desired file source instance where the
// initialization data comes from a configuration partial instance.
func (s fileSourceFactoryStrategy) CreateConfig(conf Partial) (Source, error) {
	path := conf.String("path")
	format := conf.String("format")

	return s.Create(path, format)
}
