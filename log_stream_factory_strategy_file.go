package servlet

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"strings"
)

// LogStreamFactoryStrategyFile defines a instantiation strategy to be used
// by the log stream factory for file output log stream instantiation.
type LogStreamFactoryStrategyFile struct {
	fileSystem       afero.Fs
	formatterFactory *LogFormatterFactory
}

// NewLogStreamFactoryStrategyFile instantiate a new file stream factory
// strategy that will enable the stream factory to instantiate a new file
// stream.
func NewLogStreamFactoryStrategyFile(fileSystem afero.Fs, formatterFactory *LogFormatterFactory) (LogStreamFactoryStrategy, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("invalid nil 'fileSystem' argument")
	}
	if formatterFactory == nil {
		return nil, fmt.Errorf("invalid nil 'formatterFactory' argument")
	}

	return &LogStreamFactoryStrategyFile{
		fileSystem:       fileSystem,
		formatterFactory: formatterFactory,
	}, nil
}

// Accept will check if the file stream factory strategy can instantiate a
// stream of the requested type and with the calling parameters.
func (LogStreamFactoryStrategyFile) Accept(sourceType string, args ...interface{}) bool {
	if sourceType != LogStreamTypeFile || len(args) < 4 {
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
	case LogLevel:
	default:
		return false
	}

	return true
}

// AcceptConfig will check if the stream factory strategy can instantiate a
// stream where the data to check comes from a configuration partial instance.
func (s LogStreamFactoryStrategyFile) AcceptConfig(conf ConfigPartial) (check bool) {
	defer func() {
		if r := recover(); r != nil {
			check = false
		}
	}()

	sourceType := conf.String("type")
	path := conf.String("path")
	format := conf.String("format")
	channels := s.channels(conf.Get("channels").([]interface{}))
	level := s.level(conf.String("level"))

	return s.Accept(sourceType, path, format, channels, level)
}

// Create will instantiate the desired stream instance.
func (s LogStreamFactoryStrategyFile) Create(args ...interface{}) (stream LogStream, err error) {
	defer func() {
		if r := recover(); r != nil {
			stream = nil
			err = r.(error)
		}
	}()

	path := args[0].(string)
	format := args[1].(string)
	channels := args[2].([]string)
	level := args[3].(LogLevel)

	file, err := s.fileSystem.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	formatter, err := s.formatterFactory.Create(format)
	if err != nil {
		return nil, err
	}

	return NewLogStreamFile(file, formatter, channels, level)
}

// CreateConfig will instantiate the desired stream instance where the
// initialization data comes from a configuration instance.
func (s LogStreamFactoryStrategyFile) CreateConfig(conf ConfigPartial) (stream LogStream, err error) {
	defer func() {
		if r := recover(); r != nil {
			stream = nil
			err = r.(error)
		}
	}()

	path := conf.String("path")
	format := conf.String("format")
	channels := s.channels(conf.Get("channels").([]interface{}))
	level := s.level(conf.String("level"))

	return s.Create(path, format, channels, level)
}

func (LogStreamFactoryStrategyFile) level(level string) LogLevel {
	level = strings.ToLower(level)
	if _, ok := LogLevelMap[level]; !ok {
		panic(fmt.Errorf("unrecognized logger level : %s", level))
	}
	return LogLevelMap[level]
}

func (LogStreamFactoryStrategyFile) channels(entries []interface{}) []string {
	var channels []string
	for _, channel := range entries {
		channels = append(channels, channel.(string))
	}
	return channels
}
