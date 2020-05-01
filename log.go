package servlet

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

/// ---------------------------------------------------------------------------
/// constants
/// ---------------------------------------------------------------------------

const (
	// LogFormatterFormatJSON defines the value to be used to declare a JSON
	// log formatter format.
	LogFormatterFormatJson = "json"

	// LogStreamTypeFile defines the value to be used to declare a file
	// log stream type.
	LogStreamTypeFile = "file"

	// ContainerLoggerID defines the id to be used as the default of a
	// logger instance in the application container.
	ContainerLoggerID = "servlet.log"

	// EnvContainerLoggerID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// logger id.
	EnvContainerLoggerID = "SERVLET_CONTAINER_LOGGER_ID"

	// ContainerLogFormatterFactoryID defines the id to be used as the
	// default of a logger formatter factory instance in the application
	// container.
	ContainerLogFormatterFactoryID = "servlet.log.factory.formatter"

	// EnvContainerLogFormatterFactoryID defines the name of the environment
	// variable to be checked for a overriding value for the application
	// container logger formatter factory id.
	EnvContainerLogFormatterFactoryID = "SERVLET_CONTAINER_LOGGER_FORMATTER_FACTORY_ID"

	// ContainerLogStreamFactoryID defines the id to be used as the default
	// of a logger source factory instance in the application container.
	ContainerLogStreamFactoryID = "servlet.log.factory.stream"

	// EnvContainerLogStreamFactoryID defines the name of the environment
	// variable to be checked for a overriding value for the application
	// container logger stream factory id.
	EnvContainerLogStreamFactoryID = "SERVLET_CONTAINER_LOGGER_STREAM_FACTORY_ID"

	// ContainerLogLoaderID defines the id to be used as the default of a
	// logger loader instance in the application container.
	ContainerLogLoaderID = "servlet.log.loader"

	// EnvContainerLogLoaderID defines the name of the environment
	// variable to be checked for a overriding value for the application
	// container logger loader id.
	EnvContainerLogLoaderID = "SERVLET_CONTAINER_LOGGER_LOADER_ID"
)

/// ---------------------------------------------------------------------------
/// LogLevel
/// ---------------------------------------------------------------------------

// LogLevel identifies a value type that describes a logging level.
type LogLevel int

const (
	// FATAL defines a fatal logging level.
	FATAL LogLevel = 1 + iota
	// ERROR defines a error logging level.
	ERROR
	// WARNING defines a warning logging level.
	WARNING
	// NOTICE defines a notice logging level.
	NOTICE
	// INFO defines a info logging level.
	INFO
	// DEBUG defines a debug logging level.
	DEBUG
)

// LogLevelMap defines a relation between a human-readable string
// and a code level identifier of a logging level.
var LogLevelMap = map[string]LogLevel{
	"fatal":   FATAL,
	"error":   ERROR,
	"warning": WARNING,
	"notice":  NOTICE,
	"info":    INFO,
	"debug":   DEBUG,
}

// LogLevelNameMap defines a relation between a code level identifier of a
// logging level and human-readable string representation of that level.
var LogLevelNameMap = map[LogLevel]string{
	FATAL:   "fatal",
	ERROR:   "error",
	WARNING: "warning",
	NOTICE:  "notice",
	INFO:    "info",
	DEBUG:   "debug",
}

/// ---------------------------------------------------------------------------
/// LogFormatter
/// ---------------------------------------------------------------------------

// LogFormatter interface defines the methods of a logging formatter instance
// responsible to parse a logging request into the output string.
type LogFormatter interface {
	Format(level LogLevel, message string, fields map[string]interface{}) string
}

/// ---------------------------------------------------------------------------
/// LogFormatterFactoryStrategy
/// ---------------------------------------------------------------------------

// LogFormatterFactoryStrategy interface defines the methods of the formatter
// factory strategy that can validate creation requests and instantiation
// of particular decoder.
type LogFormatterFactoryStrategy interface {
	Accept(format string, args ...interface{}) bool
	Create(args ...interface{}) (LogFormatter, error)
}

/// ---------------------------------------------------------------------------
/// LogFormatterFactory
/// ---------------------------------------------------------------------------

type LogFormatterFactory struct {
	strategies []LogFormatterFactoryStrategy
}

// NewLogFormatterFactory instantiate a new formatter factory.
func NewLogFormatterFactory() *LogFormatterFactory {
	return &LogFormatterFactory{
		strategies: []LogFormatterFactoryStrategy{},
	}
}

// Register will register a new formatter factory strategy to be used
// on requesting to create a formatter for a defined format.
func (f *LogFormatterFactory) Register(strategy LogFormatterFactoryStrategy) error {
	if f == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if strategy == nil {
		return fmt.Errorf("invalid nil 'strategy' argument")
	}

	f.strategies = append([]LogFormatterFactoryStrategy{strategy}, f.strategies...)

	return nil
}

// Create will instantiate and return a new content formatter.
func (f LogFormatterFactory) Create(format string, args ...interface{}) (LogFormatter, error) {
	for _, s := range f.strategies {
		if s.Accept(format, args...) {
			return s.Create(args...)
		}
	}
	return nil, fmt.Errorf("unrecognized format type : %s", format)
}

/// ---------------------------------------------------------------------------
/// LogJsonFormatter
/// ---------------------------------------------------------------------------

type LogJsonFormatter struct{}

// NewJSONFormatter will instantiate a new JSON formatter that will take the
// logging entry request and create the output JSON string.
func NewLogJsonFormatter() LogFormatter {
	return &LogJsonFormatter{}
}

// Format will create the output JSON string message formatted with the content
// of the passed level, fields and message
func (f LogJsonFormatter) Format(level LogLevel, message string, fields map[string]interface{}) string {
	if fields == nil {
		fields = map[string]interface{}{}
	}

	fields["time"] = time.Now().Format("2006-01-02T15:04:05.000-0700")
	fields["level"] = strings.ToUpper(LogLevelNameMap[level])
	fields["message"] = message

	bytes, _ := json.Marshal(fields)
	return string(bytes)
}

/// ---------------------------------------------------------------------------
/// LogJsonFormatterFactoryStrategy
/// ---------------------------------------------------------------------------

type LogJsonFormatterFactoryStrategy struct{}

// NewLogJsonFormatterFactoryStrategy instantiate a new json logging output
// formatter factory strategy that will enable the formatter factory to
// instantiate a new content to json formatter.
func NewLogJsonFormatterFactoryStrategy() *LogJsonFormatterFactoryStrategy {
	return &LogJsonFormatterFactoryStrategy{}
}

// Accept will check if the formatter factory strategy can instantiate a
// formatter of the requested format.
func (LogJsonFormatterFactoryStrategy) Accept(format string, _ ...interface{}) bool {
	return format == LogFormatterFormatJson
}

// Create will instantiate the desired formatter instance.
func (LogJsonFormatterFactoryStrategy) Create(_ ...interface{}) (LogFormatter, error) {
	return NewLogJsonFormatter(), nil
}

/// ---------------------------------------------------------------------------
/// LogStream
/// ---------------------------------------------------------------------------

// LogStream interface defines the interaction methods with a logging stream.
type LogStream interface {
	Close() error
	Signal(channel string, level LogLevel, message string, fields map[string]interface{}) error
	Broadcast(level LogLevel, message string, fields map[string]interface{}) error
	HasChannel(channel string) bool
	ListChannels() []string
	AddChannel(channel string)
	RemoveChannel(channel string)
	Level() LogLevel
}

/// ---------------------------------------------------------------------------
/// LogBaseStream
/// ---------------------------------------------------------------------------

type LogBaseStream struct {
	formatter LogFormatter
	channels  []string
	level     LogLevel
}

// HasChannel will validate if the stream is listening to a specific
// logging channel.
func (s LogBaseStream) HasChannel(channel string) bool {
	i := sort.SearchStrings(s.channels, channel)
	return i < len(s.channels) && s.channels[i] == channel
}

// ListChannels retrieves the list of channels that the stream is listening.
func (s LogBaseStream) ListChannels() []string {
	return s.channels
}

// AddChannel register a channel to the list of channels that the
// stream is listening.
func (s *LogBaseStream) AddChannel(channel string) {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if !s.HasChannel(channel) {
		s.channels = append(s.channels, channel)
		sort.Strings(s.channels)
	}
}

// RemoveChannel removes a channel from the list of channels that the
// stream is listening.
func (s *LogBaseStream) RemoveChannel(channel string) {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	i := sort.SearchStrings(s.channels, channel)
	if i == len(s.channels) || s.channels[i] != channel {
		return
	}
	s.channels = append(s.channels[:i], s.channels[i+1:]...)
}

// Level retrieves the logging level filter value of the stream.
func (s LogBaseStream) Level() LogLevel {
	return s.level
}

func (s LogBaseStream) format(level LogLevel, message string, fields map[string]interface{}) string {
	if s.formatter != nil {
		message = s.formatter.Format(level, message, fields)
	}
	return message
}

/// ---------------------------------------------------------------------------
/// LogStreamFactoryStrategy
/// ---------------------------------------------------------------------------

// LogStreamFactoryStrategy interface defines the methods of the stream
// factory strategy that can validate creation requests and instantiation
// of particular type of stream.
type LogStreamFactoryStrategy interface {
	Accept(sourceType string, args ...interface{}) bool
	AcceptConfig(conf ConfigPartial) bool
	Create(args ...interface{}) (LogStream, error)
	CreateConfig(conf ConfigPartial) (LogStream, error)
}

/// ---------------------------------------------------------------------------
/// LogStreamFactory
/// ---------------------------------------------------------------------------

type LogStreamFactory struct {
	strategies []LogStreamFactoryStrategy
}

// NewLogStreamFactory instantiate a new stream factory.
func NewLogStreamFactory() *LogStreamFactory {
	return &LogStreamFactory{
		strategies: []LogStreamFactoryStrategy{},
	}
}

// Register will register a new stream factory strategy to be used
// on creation requests.
func (f *LogStreamFactory) Register(strategy LogStreamFactoryStrategy) error {
	if f == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if strategy == nil {
		return fmt.Errorf("invalid nil 'strategy' argument")
	}

	f.strategies = append([]LogStreamFactoryStrategy{strategy}, f.strategies...)

	return nil
}

// Create will instantiate and return a new config stream.
func (f LogStreamFactory) Create(sourceType string, args ...interface{}) (LogStream, error) {
	for _, s := range f.strategies {
		if s.Accept(sourceType, args...) {
			return s.Create(args...)
		}
	}
	return nil, fmt.Errorf("unrecognized stream type : %s", sourceType)
}

// CreateConfig will instantiate and return a new config stream loaded by a
// configuration instance.
func (f LogStreamFactory) CreateConfig(conf ConfigPartial) (LogStream, error) {
	for _, s := range f.strategies {
		if s.AcceptConfig(conf) {
			return s.CreateConfig(conf)
		}
	}
	return nil, fmt.Errorf("unrecognized stream config : %v", conf)
}

/// ---------------------------------------------------------------------------
/// LogFileStream
/// ---------------------------------------------------------------------------

type LogFileStream struct {
	LogBaseStream
	writer io.Writer
}

// NewLogFileStream instantiate a new file stream object that will write logging
// content into a file.
func NewLogFileStream(writer io.Writer, formatter LogFormatter, channels []string, level LogLevel) (LogStream, error) {
	if formatter == nil {
		return nil, fmt.Errorf("invalid nil 'formatter' argument")
	}
	if writer == nil {
		return nil, fmt.Errorf("invalid nil 'writer' argument")
	}

	s := &LogFileStream{
		LogBaseStream: LogBaseStream{
			formatter,
			channels,
			level},
		writer: writer}

	sort.Strings(s.channels)

	return s, nil
}

// Close will terminate the stream stored writer instance.
func (s *LogFileStream) Close() (err error) {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if s.writer != nil {
		switch s.writer.(type) {
		case io.Closer:
			err = s.writer.(io.Closer).Close()
		}
		s.writer = nil
	}
	return err
}

// Signal will process the logging signal request and store the logging request
// into the underlying file if passing the channel and level filtering.
func (s LogFileStream) Signal(channel string, level LogLevel, message string, fields map[string]interface{}) error {
	i := sort.SearchStrings(s.channels, channel)
	if i == len(s.channels) || s.channels[i] != channel {
		return nil
	}
	return s.Broadcast(level, message, fields)
}

// Broadcast will process the logging signal request and store the logging
// request into the underlying file if passing the level filtering.
func (s LogFileStream) Broadcast(level LogLevel, message string, fields map[string]interface{}) error {
	if s.level < level {
		return nil
	}

	_, err := fmt.Fprintln(s.writer, s.format(level, message, fields))
	return err
}

/// ---------------------------------------------------------------------------
/// LogFileStreamFactoryStrategy
/// ---------------------------------------------------------------------------

type LogFileStreamFactoryStrategy struct {
	fileSystem       afero.Fs
	formatterFactory *LogFormatterFactory
}

// NewLogFileStreamFactoryStrategy instantiate a new file stream factory
// strategy that will enable the stream factory to instantiate a new file
// stream.
func NewLogFileStreamFactoryStrategy(fileSystem afero.Fs, formatterFactory *LogFormatterFactory) (LogStreamFactoryStrategy, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("invalid nil 'fileSystem' argument")
	}
	if formatterFactory == nil {
		return nil, fmt.Errorf("invalid nil 'formatterFactory' argument")
	}

	return &LogFileStreamFactoryStrategy{
		fileSystem:       fileSystem,
		formatterFactory: formatterFactory,
	}, nil
}

// Accept will check if the file stream factory strategy can instantiate a
// stream of the requested type and with the calling parameters.
func (LogFileStreamFactoryStrategy) Accept(sourceType string, args ...interface{}) bool {
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
func (s LogFileStreamFactoryStrategy) AcceptConfig(conf ConfigPartial) (check bool) {
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
func (s LogFileStreamFactoryStrategy) Create(args ...interface{}) (stream LogStream, err error) {
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

	return NewLogFileStream(file, formatter, channels, level)
}

// CreateConfig will instantiate the desired stream instance where the
// initialization data comes from a configuration instance.
func (s LogFileStreamFactoryStrategy) CreateConfig(conf ConfigPartial) (stream LogStream, err error) {
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

func (LogFileStreamFactoryStrategy) level(level string) LogLevel {
	level = strings.ToLower(level)
	if _, ok := LogLevelMap[level]; !ok {
		panic(fmt.Errorf("unrecognized logger level : %s", level))
	}
	return LogLevelMap[level]
}

func (LogFileStreamFactoryStrategy) channels(entries []interface{}) []string {
	var channels []string
	for _, channel := range entries {
		channels = append(channels, channel.(string))
	}
	return channels
}

/// ---------------------------------------------------------------------------
/// Log
/// ---------------------------------------------------------------------------

type Log struct {
	streams map[string]LogStream
}

// NewLog create a new logger instance.
func NewLog() *Log {
	return &Log{
		streams: map[string]LogStream{},
	}
}

// Close will terminate all the logging stream associated to the logger.
func (l *Log) Close() error {
	if l == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	for id, stream := range l.streams {
		_ = stream.Close()
		delete(l.streams, id)
	}
	return nil
}

// Signal will propagate the channel filtered logging request
// to all stored logging streams.
func (l Log) Signal(channel string, level LogLevel, message string, fields map[string]interface{}) error {
	for _, stream := range l.streams {
		if err := stream.Signal(channel, level, message, fields); err != nil {
			return err
		}
	}
	return nil
}

// Broadcast will propagate the logging request to all stored logging streams.
func (l Log) Broadcast(level LogLevel, message string, fields map[string]interface{}) error {
	for _, stream := range l.streams {
		if err := stream.Broadcast(level, message, fields); err != nil {
			return err
		}
	}
	return nil
}

// HasStream check if a stream is registered with the requested id.
func (l Log) HasStream(id string) bool {
	_, ok := l.streams[id]
	return ok
}

// AddStream registers a new stream into the logger instance.
func (l *Log) AddStream(id string, stream LogStream) error {
	if l == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if stream == nil {
		return fmt.Errorf("invalid nil 'stream' argument")
	}

	if l.HasStream(id) {
		return fmt.Errorf("duplicate id : %s", id)
	}

	l.streams[id] = stream
	return nil
}

// RemoveStream will remove a registered stream with the requested id
// from the logger.
func (l *Log) RemoveStream(id string) {
	if l == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if stream, ok := l.streams[id]; ok {
		_ = stream.Close()
		delete(l.streams, id)
	}
}

// Stream retrieve a stream from the logger that is registered with the
// requested id.
func (l Log) Stream(id string) LogStream {
	if stream, ok := l.streams[id]; ok {
		return stream
	}
	return nil
}

/// ---------------------------------------------------------------------------
/// LogLoader
/// ---------------------------------------------------------------------------

type LogLoader struct {
	logger        *Log
	streamFactory *LogStreamFactory
}

// NewLogLoader create a new logging configuration loader instance.
func NewLogLoader(logger *Log, streamFactory *LogStreamFactory) (*LogLoader, error) {
	if logger == nil {
		return nil, fmt.Errorf("invalid nil 'logger' argument")
	}
	if streamFactory == nil {
		return nil, fmt.Errorf("invalid nil 'streamFactory' argument")
	}

	return &LogLoader{
		logger:        logger,
		streamFactory: streamFactory,
	}, nil
}

// Load will parse the configuration and instantiates logging streams
// depending the data on the configuration.
func (l LogLoader) Load(c *Config) (err error) {
	if c == nil {
		return fmt.Errorf("invalid nil 'config' argument")
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	entries := c.Get("log.streams")
	if entries == nil {
		return nil
	}

	for _, entry := range entries.([]interface{}) {
		if err = l.load(entry.(ConfigPartial)); err != nil {
			return err
		}
	}

	return nil
}

func (l LogLoader) load(conf ConfigPartial) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	id := conf.String("id")

	var stream LogStream
	if stream, err = l.streamFactory.CreateConfig(conf); err != nil {
		return err
	}

	if err = l.logger.AddStream(id, stream); err != nil {
		return err
	}

	return nil
}

/// ---------------------------------------------------------------------------
/// LogParams
/// ---------------------------------------------------------------------------

// LogParams defines the logging provider parameters storing structure
// that will be needed when instantiating a new provider
type LogParams struct {
	LoggerID           string
	FileSystemID       string
	ConfigID           string
	FormatterFactoryID string
	StreamFactoryID    string
	LoaderID           string
}

// NewLogParams will instantiate a new log provider parameters
// storing instance with the servlet default values.
func NewLogParams() *LogParams {
	loggerID := ContainerLoggerID
	if env := os.Getenv(EnvContainerLoggerID); env != "" {
		loggerID = env
	}

	fileSystemID := ContainerFileSystemID
	if env := os.Getenv(EnvContainerFileSystemID); env != "" {
		fileSystemID = env
	}

	configID := ContainerConfigID
	if env := os.Getenv(EnvContainerConfigID); env != "" {
		configID = env
	}

	formatterFactoryID := ContainerLogFormatterFactoryID
	if env := os.Getenv(EnvContainerLogFormatterFactoryID); env != "" {
		formatterFactoryID = env
	}

	streamFactoryID := ContainerLogStreamFactoryID
	if env := os.Getenv(EnvContainerLogStreamFactoryID); env != "" {
		streamFactoryID = env
	}

	loaderID := ContainerLogLoaderID
	if env := os.Getenv(EnvContainerLogLoaderID); env != "" {
		loaderID = env
	}

	return &LogParams{
		LoggerID:           loggerID,
		FileSystemID:       fileSystemID,
		ConfigID:           configID,
		FormatterFactoryID: formatterFactoryID,
		StreamFactoryID:    streamFactoryID,
		LoaderID:           loaderID,
	}
}

/// ---------------------------------------------------------------------------
/// LogProvider
/// ---------------------------------------------------------------------------

type LogProvider struct {
	params *LogParams
}

// NewLogProvider will create a new logger provider instance.
func NewLogProvider(params *LogParams) *LogProvider {
	if params == nil {
		params = NewLogParams()
	}

	return &LogProvider{
		params: params,
	}
}

// Register will register the logger package instances in the
// application container.
func (p LogProvider) Register(container *AppContainer) error {
	if container == nil {
		return fmt.Errorf("invalid nil 'container' argument")
	}

	_ = container.Add(p.params.FormatterFactoryID, func(container *AppContainer) (interface{}, error) {
		formatterFactory := NewLogFormatterFactory()
		_ = formatterFactory.Register(NewLogJsonFormatterFactoryStrategy())

		return formatterFactory, nil
	})

	_ = container.Add(p.params.StreamFactoryID, func(container *AppContainer) (obj interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()

		fileSystem, err := container.Get(p.params.FileSystemID)
		if err != nil {
			return nil, err
		}

		formatterFactory, err := container.Get(p.params.FormatterFactoryID)
		if err != nil {
			return nil, err
		}

		fileStreamFactoryStrategy, _ := NewLogFileStreamFactoryStrategy(fileSystem.(afero.Fs), formatterFactory.(*LogFormatterFactory))

		streamFactory := NewLogStreamFactory()
		_ = streamFactory.Register(fileStreamFactoryStrategy)

		return streamFactory, nil
	})

	_ = container.Add(p.params.LoggerID, func(container *AppContainer) (interface{}, error) {
		return NewLog(), nil
	})

	_ = container.Add(p.params.LoaderID, func(container *AppContainer) (obj interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()

		logger, err := container.Get(p.params.LoggerID)
		if err != nil {
			return nil, err
		}

		streamFactory, err := container.Get(p.params.StreamFactoryID)
		if err != nil {
			return nil, err
		}

		return NewLogLoader(logger.(*Log), streamFactory.(*LogStreamFactory))
	})

	return nil
}

// Boot will start the logger package config instance by calling the
// logger loader with the defined provider base entry information.
func (p LogProvider) Boot(container *AppContainer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	loader, err := container.Get(p.params.LoaderID)
	if err != nil {
		return err
	}

	config, err := container.Get(p.params.ConfigID)
	if err != nil {
		return err
	}

	return loader.(*LogLoader).Load(config.(*Config))
}
