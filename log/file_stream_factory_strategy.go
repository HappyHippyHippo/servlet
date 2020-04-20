package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/happyhippyhippo/servlet/config"
	"github.com/spf13/afero"
)

const (
	// StreamTypeFile defines the value to be used to declare a file
	// log stream type.
	StreamTypeFile = "file"
)

type fileStreamFactoryStrategy struct {
	fileSystem       afero.Fs
	formatterFactory FormatterFactory
}

// NewFileStreamFactoryStrategy instantiate a new file stream factory
// strategy that will enable the stream factory to instantiate a new file
// stream.
func NewFileStreamFactoryStrategy(fileSystem afero.Fs, formatterFactory FormatterFactory) (StreamFactoryStrategy, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("Invalid nil 'fileSystem' argument")
	}
	if formatterFactory == nil {
		return nil, fmt.Errorf("Invalid nil 'formatterFactory' argument")
	}

	return &fileStreamFactoryStrategy{
		fileSystem:       fileSystem,
		formatterFactory: formatterFactory,
	}, nil
}

// Accept will check if the file stream factory strategy can instantiate a
// stream of the requested type and with the calling parameters.
func (fileStreamFactoryStrategy) Accept(stype string, args ...interface{}) bool {
	if stype != StreamTypeFile || len(args) < 4 {
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

	switch args[2].(type) {
	case []string:
	default:
		return false
	}

	switch args[3].(type) {
	case Level:
	default:
		return false
	}

	return true
}

// AcceptConfig will check if the stream factory strategy can instantiate a
// stream where the data to check comes from a configuration partial instance.
func (s fileStreamFactoryStrategy) AcceptConfig(conf config.Partial) (check bool) {
	defer func() {
		if r := recover(); r != nil {
			check = false
		}
	}()

	stype := conf.String("type")
	path := conf.String("path")
	format := conf.String("format")
	channels := s.channels(conf.Get("channels").([]interface{}))
	level := s.level(conf.String("level"))

	return s.Accept(stype, path, format, channels, level)
}

// Create will instantiate the desired stream instance.
func (s fileStreamFactoryStrategy) Create(args ...interface{}) (Stream, error) {
	path := args[0].(string)
	format := args[1].(string)
	channels := args[2].([]string)
	level := args[3].(Level)

	file, err := s.fileSystem.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	formatter, err := s.formatterFactory.Create(format)
	if err != nil {
		return nil, err
	}

	return NewFileStream(file, formatter, channels, level)
}

// CreateConfig will instantiate the desired stream instance where the
// initialization data comes from a configuration instance.
func (s fileStreamFactoryStrategy) CreateConfig(conf config.Partial) (Stream, error) {
	path := conf.String("path")
	format := conf.String("format")
	channels := s.channels(conf.Get("channels").([]interface{}))
	level := s.level(conf.String("level"))

	return s.Create(path, format, channels, level)
}

func (fileStreamFactoryStrategy) level(level string) Level {
	level = strings.ToLower(level)
	if _, ok := LevelMap[level]; !ok {
		panic(fmt.Errorf("Unrecognized logger level : %s", level))
	}
	return LevelMap[level]
}

func (fileStreamFactoryStrategy) channels(entries []interface{}) []string {
	var channels = []string{}
	for _, channel := range entries {
		channels = append(channels, channel.(string))
	}
	return channels
}
