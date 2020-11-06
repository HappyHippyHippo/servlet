package servlet

import (
	"fmt"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

/// ---------------------------------------------------------------------------
/// constants
/// ---------------------------------------------------------------------------

const (
	// ConfigDecoderFormatYAML defines the value to be used to declare a YAML
	// config source format.
	ConfigDecoderFormatYAML = "yaml"

	// ConfigSourceTypeFile defines the value to be used to declare a
	// simple file config source type.
	ConfigSourceTypeFile = "file"

	// ConfigSourceTypeObservableFile defines the value to be used to declare a
	// observable file config source type.
	ConfigSourceTypeObservableFile = "observable_file"

	// ConfigSourceTypeEnv defines the value to be used to declare a
	// environment config source type.
	ConfigSourceTypeEnv = "env"

	// ContainerConfigID defines the id to be used as the default of a
	// config instance in the application container.
	ContainerConfigID = "servlet.config"

	// EnvContainerConfigID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// config id.
	EnvContainerConfigID = "SERVLET_CONTAINER_CONFIG_ID"

	// ContainerConfigDecoderFactoryID defines the id to be used as the default of a
	// config decoder factory instance in the application container.
	ContainerConfigDecoderFactoryID = "servlet.config.factory.decoder"

	// EnvContainerConfigDecoderFactoryID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// config decoder factory id.
	EnvContainerConfigDecoderFactoryID = "SERVLET_CONTAINER_CONFIG_DECODER_FACTORY_ID"

	// ContainerConfigSourceFactoryID defines the id to be used as the default of a
	// config source factory instance in the application container.
	ContainerConfigSourceFactoryID = "servlet.config.factory.source"

	// EnvContainerConfigSourceFactoryID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// config source factory id.
	EnvContainerConfigSourceFactoryID = "SERVLET_CONTAINER_CONFIG_SOURCE_FACTORY_ID"

	// ContainerConfigLoaderID defines the id to be used as the default of a
	// config loader instance in the application container.
	ContainerConfigLoaderID = "servlet.config.loader"

	// EnvContainerConfigLoaderID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// config loaded id.
	EnvContainerConfigLoaderID = "SERVLET_CONTAINER_CONFIG_LOADER_ID"

	// ConfigObserveFrequency defines the id to be used as the default of a
	// config observable source frequency time.
	ConfigObserveFrequency = time.Second * 0

	// EnvConfigObserveFrequency defines the name of the environment variable
	// to be checked for a overriding value for the config observe frequency.
	EnvConfigObserveFrequency = "SERVLET_CONFIG_OBSERVE_FREQUENCY"

	// ConfigBaseSourceActive defines the base config source active flag
	// used to signal the config loader to load the base source or not
	ConfigBaseSourceActive = true

	// EnvConfigBaseSourceActive defines the name of the environment variable
	// to be checked for a overriding value for the config base source active.
	EnvConfigBaseSourceActive = "SERVLET_CONFIG_BASE_SOURCE_ACTIVE"

	// ConfigBaseSourceID defines the id to be used as the default of the
	// base config source id to be used as the loader entry.
	ConfigBaseSourceID = "base"

	// EnvConfigBaseSourceID defines the name of the environment variable
	// to be checked for a overriding value for the config base source id.
	EnvConfigBaseSourceID = "SERVLET_CONFIG_BASE_SOURCE_ID"

	// ConfigBaseSourcePath defines the base config source path
	// to be used as the loader entry.
	ConfigBaseSourcePath = "config/config.yaml"

	// EnvConfigBaseSourcePath defines the name of the environment variable
	// to be checked for a overriding value for the config base source path.
	EnvConfigBaseSourcePath = "SERVLET_CONFIG_BASE_SOURCE_PATH"

	// ConfigBaseSourceFormat defines the base config source format
	// to be used as the loader entry.
	ConfigBaseSourceFormat = ConfigDecoderFormatYAML

	// EnvConfigBaseSourceFormat defines the name of the environment variable
	// to be checked for a overriding value for the config base source format.
	EnvConfigBaseSourceFormat = "SERVLET_CONFIG_BASE_SOURCE_FORMAT"
)

/// ---------------------------------------------------------------------------
/// ConfigPartial
/// ---------------------------------------------------------------------------

// ConfigPartial defined a type used to store configuration information.
type ConfigPartial map[interface{}]interface{}

// Has will check if a requested path exists in the config partial.
func (p ConfigPartial) Has(path string) bool {
	it := p
	nodes := strings.Split(path, ".")
	for i, node := range nodes {
		if node == "" {
			continue
		}

		switch it[node].(type) {
		case ConfigPartial:
			it = it[node].(ConfigPartial)
		case nil:
			return false
		default:
			return i == len(nodes)-1
		}
	}

	return true
}

// Get will retrieve the value stored in the requested path.
// If the path does not exists, then the value nil will be returned. Or, if
// a default value was given as the optional extra argument, then it will
// be returned instead of the standard nil value.
func (p ConfigPartial) Get(path string, def ...interface{}) interface{} {
	it := p
	nodes := strings.Split(path, ".")
	for i, node := range nodes {
		if node == "" {
			continue
		}

		if _, ok := it[node]; !ok {
			if len(def) > 0 {
				return def[0]
			}
			return nil
		}

		switch it[node].(type) {
		case ConfigPartial:
			it = it[node].(ConfigPartial)
		case nil:
			return nil
		default:
			if i != len(nodes)-1 {
				if len(def) > 0 {
					return def[0]
				}
				return nil
			}
			return it[node]
		}
	}

	return it
}

// Int will return the casting to int of the stored value in the
// requested path. If the value retrieved was not found or returned nil, then
// the default optional argument will be returned if given.
func (p ConfigPartial) Int(path string, def ...int) int {
	value := p.Get(path)
	if value == nil && len(def) > 0 {
		return def[0]
	}
	return value.(int)
}

// String will return the casting to string of the stored value in the
// requested path. If the value retrieved was not found or returned nil, then
// the default optional argument will be returned if given.
func (p ConfigPartial) String(path string, def ...string) string {
	value := p.Get(path)
	if value == nil && len(def) > 0 {
		return def[0]
	}
	return p.Get(path).(string)
}

// Config will return the casting to a config partial of the stored
// value in the requested path. If the value retrieved was not found or
// returned nil, then the default optional argument will be returned if given.
func (p ConfigPartial) Config(path string, def ...ConfigPartial) ConfigPartial {
	value := p.Get(path)
	if value == nil && len(def) > 0 {
		return def[0]
	}
	return p.Get(path).(ConfigPartial)
}

func (p ConfigPartial) merge(p2 ConfigPartial) ConfigPartial {
	for key, value := range p2 {
		switch value.(type) {
		case ConfigPartial:
			switch p[key].(type) {
			case ConfigPartial:
				p[key] = p[key].(ConfigPartial).merge(value.(ConfigPartial))
			default:
				p[key] = value
			}
		default:
			p[key] = value
		}
	}
	return p
}

/// ---------------------------------------------------------------------------
/// ConfigDecoder
/// ---------------------------------------------------------------------------

// ConfigDecoder interface defines the interaction methods to a config content
// decoder used to parse the source content into a application usable
// configuration partial instance.
type ConfigDecoder interface {
	Close()
	Decode() (ConfigPartial, error)
}

/// ---------------------------------------------------------------------------
/// ConfigDecoderFactoryStrategy
/// ---------------------------------------------------------------------------

// ConfigDecoderFactoryStrategy interface defines the methods of the decoder
// factory strategy that can validate creation requests and instantiation of a
// particular decoder.
type ConfigDecoderFactoryStrategy interface {
	Accept(format string, args ...interface{}) bool
	Create(args ...interface{}) (ConfigDecoder, error)
}

/// ---------------------------------------------------------------------------
/// ConfigDecoderFactory
/// ---------------------------------------------------------------------------

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

/// ---------------------------------------------------------------------------
/// ConfigYamlDecoder
/// ---------------------------------------------------------------------------

type underlyingConfigYamlDecoder interface {
	Decode(partial interface{}) error
}

// ConfigYamlDecoder defines an instance used to decode s YAML encoded config
// source stream
type ConfigYamlDecoder struct {
	reader  io.Reader
	decoder underlyingConfigYamlDecoder
}

// NewConfigYamlDecoder instantiate a new yaml configuration decoder object
// used to parse a yaml configuration source into a config partial.
func NewConfigYamlDecoder(reader io.Reader) (*ConfigYamlDecoder, error) {
	if reader == nil {
		return nil, fmt.Errorf("invalid nil 'reader' argument")
	}

	return &ConfigYamlDecoder{
		reader:  reader,
		decoder: yaml.NewDecoder(reader),
	}, nil
}

// Close terminate the decoder, closing the associated reader.
func (d *ConfigYamlDecoder) Close() {
	if d == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if d.reader != nil {
		switch d.reader.(type) {
		case io.Closer:
			_ = d.reader.(io.Closer).Close()
		}
		d.reader = nil
	}
}

// Decode parse the associated configuration source reader content
// into a configuration partial.
func (d ConfigYamlDecoder) Decode() (ConfigPartial, error) {
	p := ConfigPartial{}
	if err := d.decoder.Decode(&p); err != nil {
		return nil, err
	}
	return p, nil
}

/// ---------------------------------------------------------------------------
/// ConfigYamlDecoderFactoryStrategy
/// ---------------------------------------------------------------------------

// ConfigYamlDecoderFactoryStrategy defines a strategy used to instantiate
// a YAML config stream decoder.
type ConfigYamlDecoderFactoryStrategy struct{}

// NewConfigYamlDecoderFactoryStrategy instantiate a new yaml decoder factory
// strategy that will enable the decoder factory to instantiate a new yaml
// decoder.
func NewConfigYamlDecoderFactoryStrategy() *ConfigYamlDecoderFactoryStrategy {
	return &ConfigYamlDecoderFactoryStrategy{}
}

// Accept will check if the decoder factory strategy can instantiate a
// decoder giving the format and the creation request parameters.
func (ConfigYamlDecoderFactoryStrategy) Accept(format string, args ...interface{}) bool {
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
func (ConfigYamlDecoderFactoryStrategy) Create(args ...interface{}) (ConfigDecoder, error) {
	reader := args[0].(io.Reader)

	return NewConfigYamlDecoder(reader)
}

/// ---------------------------------------------------------------------------
/// ConfigSource
/// ---------------------------------------------------------------------------

// ConfigSource defines the base interface of a config source.
type ConfigSource interface {
	Close()
	Has(path string) bool
	Get(path string) interface{}
}

/// ---------------------------------------------------------------------------
/// ConfigBaseSource
/// ---------------------------------------------------------------------------

// ConfigBaseSource defines a base code of a config source instance.
type ConfigBaseSource struct {
	mutex   sync.Locker
	partial ConfigPartial
}

// Close method used to be compliant with the container Closable interface.
func (*ConfigBaseSource) Close() {}

// Has will check if the requested path is present in the source
// configuration content.
func (s *ConfigBaseSource) Has(path string) bool {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.partial.Has(path)
}

// Get will retrieve the value stored in the requested path present in the
// configuration content.
// If the path does not exists, then the value nil will be returned.
// This method will mostly be used by the config object to obtain the full
// content of the source to aggregate all the data into his internal storing
// partial instance.
func (s *ConfigBaseSource) Get(path string) interface{} {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.partial.Get(path)
}

/// ---------------------------------------------------------------------------
/// ConfigObservableSource
/// ---------------------------------------------------------------------------

// ConfigObservableSource interface extends the Source interface with methods
// specific to sources that will be checked for updates in a regular
// periodicity defined in the config object where the source will be
// registered.
type ConfigObservableSource interface {
	ConfigSource
	Reload() (bool, error)
}

/// ---------------------------------------------------------------------------
/// ConfigSourceFactoryStrategy
/// ---------------------------------------------------------------------------

// ConfigSourceFactoryStrategy interface defines the methods of the source
// factory strategy that will be used instantiate a particular source type.
type ConfigSourceFactoryStrategy interface {
	Accept(sourceType string, args ...interface{}) bool
	AcceptConfig(conf ConfigPartial) bool
	Create(args ...interface{}) (ConfigSource, error)
	CreateConfig(conf ConfigPartial) (ConfigSource, error)
}

/// ---------------------------------------------------------------------------
/// ConfigSourceFactory
/// ---------------------------------------------------------------------------

// ConfigSourceFactory defines a config source factory that uses a list of
// registered instantiation strategies to perform the config source
// instantiation.
type ConfigSourceFactory struct {
	strategies []ConfigSourceFactoryStrategy
}

// NewConfigSourceFactory instantiate a new source factory.
func NewConfigSourceFactory() *ConfigSourceFactory {
	return &ConfigSourceFactory{
		strategies: []ConfigSourceFactoryStrategy{},
	}
}

// Register will register a new source factory strategy to be used
// on creation request.
func (f *ConfigSourceFactory) Register(strategy ConfigSourceFactoryStrategy) error {
	if f == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if strategy == nil {
		return fmt.Errorf("invalid nil 'strategy' argument")
	}

	f.strategies = append([]ConfigSourceFactoryStrategy{strategy}, f.strategies...)

	return nil
}

// Create will instantiate and return a new config source by the type requested.
func (f ConfigSourceFactory) Create(sourceType string, args ...interface{}) (ConfigSource, error) {
	for _, s := range f.strategies {
		if s.Accept(sourceType, args...) {
			return s.Create(args...)
		}
	}
	return nil, fmt.Errorf("unrecognized source type : %s", sourceType)
}

// CreateConfig will instantiate and return a new config source where the
// data used to decide the strategy to be used and also the initialization
// data comes from a configuration storing partial instance.
func (f ConfigSourceFactory) CreateConfig(conf ConfigPartial) (ConfigSource, error) {
	for _, s := range f.strategies {
		if s.AcceptConfig(conf) {
			return s.CreateConfig(conf)
		}
	}
	return nil, fmt.Errorf("unrecognized source config : %v", conf)
}

/// ---------------------------------------------------------------------------
/// ConfigFileSource
/// ---------------------------------------------------------------------------

// ConfigFileSource defines an instance of a file stream configuration source.
type ConfigFileSource struct {
	ConfigBaseSource
	path           string
	format         string
	fileSystem     afero.Fs
	decoderFactory *ConfigDecoderFactory
}

// NewConfigFileSource instantiate a new source that treats a file as
// the origin of the configuration content.
func NewConfigFileSource(path string, format string, fileSystem afero.Fs, decoderFactory *ConfigDecoderFactory) (*ConfigFileSource, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("invalid nil 'fileSystem' argument")
	}
	if decoderFactory == nil {
		return nil, fmt.Errorf("invalid nil 'decoderFactory' argument")
	}

	s := &ConfigFileSource{
		ConfigBaseSource: ConfigBaseSource{
			mutex:   &sync.Mutex{},
			partial: nil,
		},
		path:           path,
		format:         format,
		fileSystem:     fileSystem,
		decoderFactory: decoderFactory,
	}

	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *ConfigFileSource) load() error {
	file, err := s.fileSystem.OpenFile(s.path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	decoder, err := s.decoderFactory.Create(s.format, file)
	if err != nil {
		_ = file.Close()
		return err
	}
	defer decoder.Close()

	partial, err := decoder.Decode()
	if err != nil {
		return err
	}

	s.mutex.Lock()
	s.partial = partial
	s.mutex.Unlock()

	return nil
}

/// ---------------------------------------------------------------------------
/// ConfigFileSourceFactoryStrategy
/// ---------------------------------------------------------------------------

// ConfigFileSourceFactoryStrategy defines a config file source instantiation
// strategy to be used by the config sources factory instance.
type ConfigFileSourceFactoryStrategy struct {
	fileSystem     afero.Fs
	decoderFactory *ConfigDecoderFactory
}

// NewConfigFileSourceFactoryStrategy instantiate a new file source factory
// strategy that will enable the source factory to instantiate a new
// file configuration source.
func NewConfigFileSourceFactoryStrategy(fileSystem afero.Fs, decoderFactory *ConfigDecoderFactory) (*ConfigFileSourceFactoryStrategy, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("invalid nil 'fileSystem' argument")
	}
	if decoderFactory == nil {
		return nil, fmt.Errorf("invalid nil 'decoderFactory' argument")
	}

	return &ConfigFileSourceFactoryStrategy{
		fileSystem:     fileSystem,
		decoderFactory: decoderFactory,
	}, nil
}

// Accept will check if the source factory strategy can instantiate a
// new source of the requested type. Also, validates that there is the path
// and content format extra parameters, and thar this parameters are strings.
func (ConfigFileSourceFactoryStrategy) Accept(sourceType string, args ...interface{}) bool {
	if sourceType != ConfigSourceTypeFile || len(args) < 2 {
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
func (s ConfigFileSourceFactoryStrategy) AcceptConfig(conf ConfigPartial) (check bool) {
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

// Create will instantiate the desired file source instance.
func (s ConfigFileSourceFactoryStrategy) Create(args ...interface{}) (source ConfigSource, err error) {
	defer func() {
		if r := recover(); r != nil {
			source = nil
			err = r.(error)
		}
	}()

	path := args[0].(string)
	format := args[1].(string)

	return NewConfigFileSource(path, format, s.fileSystem, s.decoderFactory)
}

// CreateConfig will instantiate the desired file source instance where the
// initialization data comes from a configuration partial instance.
func (s ConfigFileSourceFactoryStrategy) CreateConfig(conf ConfigPartial) (source ConfigSource, err error) {
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

/// ---------------------------------------------------------------------------
/// ConfigObservableFileSource
/// ---------------------------------------------------------------------------

// ConfigObservableFileSource defines an instance of a file stream
// configuration source that will be checked for changes periodically in a
// config defined frequency.
type ConfigObservableFileSource struct {
	ConfigFileSource
	timestamp time.Time
}

// NewConfigObservableFileSource instantiate a new source that treats a file
// as the origin of the configuration content. This file source will be
// periodically checked for changes and loaded if so.
func NewConfigObservableFileSource(path string, format string, fileSystem afero.Fs, decoderFactory *ConfigDecoderFactory) (*ConfigObservableFileSource, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("invalid nil 'fileSystem' argument")
	}
	if decoderFactory == nil {
		return nil, fmt.Errorf("invalid nil 'decoderFactory' argument")
	}

	s := &ConfigObservableFileSource{
		ConfigFileSource: ConfigFileSource{
			ConfigBaseSource: ConfigBaseSource{
				mutex:   &sync.RWMutex{},
				partial: nil,
			},
			path:           path,
			format:         format,
			fileSystem:     fileSystem,
			decoderFactory: decoderFactory,
		},
		timestamp: time.Unix(0, 0),
	}

	if _, err := s.Reload(); err != nil {
		return nil, err
	}
	return s, nil
}

// Reload will check if the source has been updated, and, if so, reload the
// source configuration partial content.
func (s *ConfigObservableFileSource) Reload() (bool, error) {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	fileInfo, err := s.fileSystem.Stat(s.path)
	if err != nil {
		return false, err
	}

	info := fileInfo.ModTime()
	if s.timestamp.Equal(time.Unix(0, 0)) || s.timestamp.Before(info) {
		if err := s.load(); err != nil {
			return false, err
		}
		s.mutex.Lock()
		s.timestamp = info
		s.mutex.Unlock()
		return true, nil
	}
	return false, nil
}

/// ---------------------------------------------------------------------------
/// ConfigObservableFileSourceFactoryStrategy
/// ---------------------------------------------------------------------------

// ConfigObservableFileSourceFactoryStrategy defines a observable config file
// source instantiation strategy to be used by the config sources factory
// instance.
type ConfigObservableFileSourceFactoryStrategy struct {
	fileSystem     afero.Fs
	decoderFactory *ConfigDecoderFactory
}

// NewConfigObservableFileSourceFactoryStrategy instantiate a new observable file
// source factory strategy that will enable the source factory to instantiate
// a new observable file configuration source.
func NewConfigObservableFileSourceFactoryStrategy(fileSystem afero.Fs, decoderFactory *ConfigDecoderFactory) (*ConfigObservableFileSourceFactoryStrategy, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("invalid nil 'fileSystem' argument")
	}
	if decoderFactory == nil {
		return nil, fmt.Errorf("invalid nil 'decoderFactory' argument")
	}

	return &ConfigObservableFileSourceFactoryStrategy{
		fileSystem:     fileSystem,
		decoderFactory: decoderFactory,
	}, nil
}

// Accept will check if the source factory strategy can instantiate a
// new source of the requested type. Also, validates that there is the path
// and content format extra parameters, and thar this parameters are strings.
func (ConfigObservableFileSourceFactoryStrategy) Accept(sourceType string, args ...interface{}) bool {
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
func (s ConfigObservableFileSourceFactoryStrategy) AcceptConfig(conf ConfigPartial) (check bool) {
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
func (s ConfigObservableFileSourceFactoryStrategy) Create(args ...interface{}) (source ConfigSource, err error) {
	defer func() {
		if r := recover(); r != nil {
			source = nil
			err = r.(error)
		}
	}()

	path := args[0].(string)
	format := args[1].(string)

	return NewConfigObservableFileSource(path, format, s.fileSystem, s.decoderFactory)
}

// CreateConfig will instantiate the desired observable file source instance
// where the initialization data comes from a configuration partial instance.
func (s ConfigObservableFileSourceFactoryStrategy) CreateConfig(conf ConfigPartial) (source ConfigSource, err error) {
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

/// ---------------------------------------------------------------------------
/// ConfigEnvSource
/// ---------------------------------------------------------------------------

// ConfigEnvSource defines an instance of a environment variables stream
// configuration source.
type ConfigEnvSource struct {
	ConfigBaseSource
	mapper map[string]string
}

// NewConfigEnvSource instantiate a new source that read a list of environment
// variables into mapped config paths.
func NewConfigEnvSource(mapper map[string]string) (*ConfigEnvSource, error) {
	s := &ConfigEnvSource{
		ConfigBaseSource: ConfigBaseSource{
			mutex:   &sync.Mutex{},
			partial: ConfigPartial{},
		},
		mapper: mapper,
	}

	_ = s.load()

	return s, nil
}

func (s *ConfigEnvSource) load() error {
	for v, p := range s.mapper {
		if env := os.Getenv(v); env != "" {
			step := s.partial
			sections := strings.Split(p, ".")
			for i, section := range sections {
				if i != len(sections)-1 {
					if _, ok := step[section]; ok == false {
						step[section] = ConfigPartial{}
					}

					switch step[section].(type) {
					case ConfigPartial:
					default:
						step[section] = ConfigPartial{}
					}

					step = step[section].(ConfigPartial)
				} else {
					step[section] = env
				}
			}
		}
	}

	return nil
}

/// ---------------------------------------------------------------------------
/// ConfigEnvSourceFactoryStrategy
/// ---------------------------------------------------------------------------

// ConfigEnvSourceFactoryStrategy defines a environment config source
// instantiation strategy to be used by the config sources factory
// instance.
type ConfigEnvSourceFactoryStrategy struct{}

// NewConfigEnvSourceFactoryStrategy instantiate a new environment
// source factory strategy that will enable the source factory to instantiate
// a new observable file configuration source.
func NewConfigEnvSourceFactoryStrategy() (*ConfigEnvSourceFactoryStrategy, error) {
	return &ConfigEnvSourceFactoryStrategy{}, nil
}

// Accept will check if the source factory strategy can instantiate a
// new source of the requested type. Also, validates that there is the path
// and content format extra parameters, and thar this parameters are strings.
func (ConfigEnvSourceFactoryStrategy) Accept(sourceType string, args ...interface{}) bool {
	if sourceType != ConfigSourceTypeEnv || len(args) < 1 {
		return false
	}

	switch args[0].(type) {
	case map[string]string:
	default:
		return false
	}

	return true
}

// AcceptConfig will check if the source factory strategy can instantiate a
// source where the data to check comes from a configuration partial instance.
func (s ConfigEnvSourceFactoryStrategy) AcceptConfig(conf ConfigPartial) (check bool) {
	defer func() {
		if r := recover(); r != nil {
			check = false
		}
	}()

	sourceType := conf.String("type")
	mapping := conf.Get("mapping")

	return s.Accept(sourceType, mapping)
}

// Create will instantiate the desired environment source instance.
func (s ConfigEnvSourceFactoryStrategy) Create(args ...interface{}) (source ConfigSource, err error) {
	defer func() {
		if r := recover(); r != nil {
			source = nil
			err = r.(error)
		}
	}()

	mappings := args[0].(map[string]string)

	return NewConfigEnvSource(mappings)
}

// CreateConfig will instantiate the desired environment source instance
// where the initialization data comes from a configuration partial instance.
func (s ConfigEnvSourceFactoryStrategy) CreateConfig(conf ConfigPartial) (source ConfigSource, err error) {
	defer func() {
		if r := recover(); r != nil {
			source = nil
			err = r.(error)
		}
	}()

	mapping := map[string]string{}
	for k, v := range conf.Get("mapping").(ConfigPartial) {
		mapping[k.(string)] = v.(string)
	}

	return s.Create(mapping)
}

/// ---------------------------------------------------------------------------
/// ConfigObserver
/// ---------------------------------------------------------------------------

// ConfigObserver callback function used to be called when a observed
// configuration path has changed.
type ConfigObserver func(interface{}, interface{})

/// ---------------------------------------------------------------------------
/// Config
/// ---------------------------------------------------------------------------

type configRefSource struct {
	id       string
	priority int
	source   ConfigSource
}

type configRefSourceSortByPriority []configRefSource

func (sources configRefSourceSortByPriority) Len() int {
	return len(sources)
}

func (sources configRefSourceSortByPriority) Swap(i, j int) {
	sources[i], sources[j] = sources[j], sources[i]
}

func (sources configRefSourceSortByPriority) Less(i, j int) bool {
	return sources[i].priority < sources[j].priority
}

type configRefObserver struct {
	path     string
	current  interface{}
	callback ConfigObserver
}

// Config defines the instance of a configuration managing structure.
type Config struct {
	mutex     sync.Locker
	sources   []configRefSource
	observers []configRefObserver
	partial   ConfigPartial
	loader    *RecurringTrigger
}

// NewConfig instantiate a new configuration object.
// This object will manage a series of sources, along side of the ability of
// registration of configuration path/values observer callbacks that will be
// called whenever the value has changed.
func NewConfig(period time.Duration) (*Config, error) {
	var c *Config

	var loader *RecurringTrigger
	if period != 0 {
		loader, _ = NewRecurringTrigger(period, func() error { return c.reload() })
	}

	c = &Config{
		mutex:     &sync.Mutex{},
		sources:   []configRefSource{},
		observers: []configRefObserver{},
		partial:   ConfigPartial{},
		loader:    loader,
	}

	return c, nil
}

// Close terminates the config instance.
// This will stop the observer trigger and call close on all registered sources.
func (c *Config) Close() {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if c.loader != nil {
		c.loader.Stop()
		c.loader = nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, reg := range c.sources {
		reg.source.Close()
	}
}

// Has will check if a path has been loaded.
// This means that if the values has been loaded by any registered source.
func (c *Config) Has(path string) bool {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.partial.Has(path)
}

// Get will retrieve a configuration value loaded from a source.
func (c *Config) Get(path string, def ...interface{}) interface{} {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.partial.Get(path, def...)
}

// GetBool will retrieve a configuration value loaded from a
// source as a boolean.
func (c *Config) GetBool(path string, def ...bool) bool {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(bool)
	}
	return c.Get(path).(bool)
}

// GetInt will retrieve a configuration value loaded from a source as a int.
func (c *Config) GetInt(path string, def ...int) int {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(int)
	}
	return c.Get(path).(int)
}

// GetInt8 will retrieve a configuration value loaded from a source as a int8.
func (c *Config) GetInt8(path string, def ...int8) int8 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(int8)
	}
	return c.Get(path).(int8)
}

// GetInt16 will retrieve a configuration value loaded from a source as a int16.
func (c *Config) GetInt16(path string, def ...int16) int16 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(int16)
	}
	return c.Get(path).(int16)
}

// GetInt32 will retrieve a configuration value loaded from a source as a int32.
func (c *Config) GetInt32(path string, def ...int32) int32 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(int32)
	}
	return c.Get(path).(int32)
}

// GetInt64 will retrieve a configuration value loaded from a source as a int64.
func (c *Config) GetInt64(path string, def ...int64) int64 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(int64)
	}
	return c.Get(path).(int64)
}

// GetUInt will retrieve a configuration value loaded from a source as a uint.
func (c *Config) GetUInt(path string, def ...uint) uint {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(uint)
	}
	return c.Get(path).(uint)
}

// GetUInt8 will retrieve a configuration value loaded from a source as a uint8.
func (c *Config) GetUInt8(path string, def ...uint8) uint8 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(uint8)
	}
	return c.Get(path).(uint8)
}

// GetUInt16 will retrieve a configuration value loaded from a
// source as a uint16.
func (c *Config) GetUInt16(path string, def ...uint16) uint16 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(uint16)
	}
	return c.Get(path).(uint16)
}

// GetUInt32 will retrieve a configuration value loaded from a
// source as a uint32.
func (c *Config) GetUInt32(path string, def ...uint32) uint32 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(uint32)
	}
	return c.Get(path).(uint32)
}

// GetUInt64 will retrieve a configuration value loaded from a
// source as a uint64.
func (c *Config) GetUInt64(path string, def ...uint64) uint64 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(uint64)
	}
	return c.Get(path).(uint64)
}

// GetFloat32 will retrieve a configuration value loaded from a
// source as a float32.
func (c *Config) GetFloat32(path string, def ...float32) float32 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(float32)
	}
	return c.Get(path).(float32)
}

// GetFloat64 will retrieve a configuration value loaded from a
// source as a float64.
func (c *Config) GetFloat64(path string, def ...float64) float64 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(float64)
	}
	return c.Get(path).(float64)
}

// GetComplex64 will retrieve a configuration value loaded from a
// source as a complex64.
func (c *Config) GetComplex64(path string, def ...complex64) complex64 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(complex64)
	}
	return c.Get(path).(complex64)
}

// GetComplex128 will retrieve a configuration value loaded from a
// source as a complex128.
func (c *Config) GetComplex128(path string, def ...complex128) complex128 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(complex128)
	}
	return c.Get(path).(complex128)
}

// GetRune will retrieve a configuration value loaded from a source as a rune.
func (c *Config) GetRune(path string, def ...rune) rune {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(rune)
	}
	return c.Get(path).(rune)
}

// GetString will retrieve a configuration value loaded from a
// source as a string.
func (c *Config) GetString(path string, def ...string) string {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(string)
	}
	return c.Get(path).(string)
}

// HasSource check if a source with a specific id has been registered.
func (c *Config) HasSource(id string) bool {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, reg := range c.sources {
		if reg.id == id {
			return true
		}
	}
	return false
}

// AddSource register a new source with a specific id with a given priority.
func (c *Config) AddSource(id string, priority int, source ConfigSource) error {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if source == nil {
		return fmt.Errorf("invalid nil 'source' argument")
	}
	if c.HasSource(id) {
		return fmt.Errorf("duplicate source id : %s", id)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.sources = append(c.sources, configRefSource{id, priority, source})
	sort.Sort(configRefSourceSortByPriority(c.sources))
	c.rebuild()

	return nil
}

// RemoveSource remove a source from the registration list
// of the configuration. This will also update the configuration content and
// re-validate the observed paths.
func (c *Config) RemoveSource(id string) {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, reg := range c.sources {
		if reg.id == id {
			reg.source.Close()
			c.sources = append(c.sources[:i], c.sources[i+1:]...)
			c.rebuild()
			return
		}
	}
}

// Source retrieve a previously registered source with a requested id.
func (c *Config) Source(id string) (ConfigSource, error) {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, reg := range c.sources {
		if reg.id == id {
			return reg.source, nil
		}
	}
	return nil, fmt.Errorf("source not found : %s", id)
}

// SourcePriority set a priority value of a previously registered
// source with the specified id. This may change the defined values if there
// was a override process of the configuration paths of the changing source.
func (c *Config) SourcePriority(id string, priority int) error {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, reg := range c.sources {
		if reg.id == id {
			reg.priority = priority
			sort.Sort(configRefSourceSortByPriority(c.sources))
			c.rebuild()

			return nil
		}
	}
	return fmt.Errorf("source not found : %s", id)
}

// HasObserver check if there is a observer to a configuration value path.
func (c *Config) HasObserver(path string) bool {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, reg := range c.observers {
		if reg.path == path {
			return true
		}
	}
	return false
}

// AddObserver register a new observer to a configuration path.
func (c *Config) AddObserver(path string, callback ConfigObserver) error {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if callback == nil {
		return fmt.Errorf("invalid nil 'callback' argument")
	}

	value := c.Get(path)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.observers = append(c.observers, configRefObserver{path, value, callback})

	return nil
}

// RemoveObserver remove a observer to a configuration path.
func (c *Config) RemoveObserver(path string) {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, reg := range c.observers {
		if reg.path == path {
			c.observers = append(c.observers[:i], c.observers[i+1:]...)
			return
		}
	}
}

func (c *Config) reload() error {
	rebuild := false
	for _, ref := range c.sources {
		switch s := ref.source.(type) {
		case ConfigObservableSource:
			changed, _ := s.Reload()
			rebuild = rebuild || changed
		}
	}

	if rebuild {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		c.rebuild()
	}

	return nil
}

func (c *Config) rebuild() {
	p := ConfigPartial{}
	for _, reg := range c.sources {
		p = p.merge(reg.source.Get("").(ConfigPartial))
	}

	c.partial = p

	for _, observer := range c.observers {
		updated := c.partial.Get(observer.path)
		if !reflect.DeepEqual(observer.current, updated) {
			old := observer.current
			observer.current = updated

			observer.callback(old, updated)
		}
	}
}

/// ---------------------------------------------------------------------------
/// ConfigLoader
/// ---------------------------------------------------------------------------

// ConfigLoader defines the config instantiation and initialization of a new
// config managing structure.
type ConfigLoader struct {
	config        *Config
	sourceFactory *ConfigSourceFactory
}

// NewConfigLoader instantiate a new configuration loader.
func NewConfigLoader(config *Config, sourceFactory *ConfigSourceFactory) (*ConfigLoader, error) {
	if config == nil {
		return nil, fmt.Errorf("invalid nil 'config' argument")
	}
	if sourceFactory == nil {
		return nil, fmt.Errorf("invalid nil 'sourceFactory' argument")
	}

	return &ConfigLoader{
		config:        config,
		sourceFactory: sourceFactory,
	}, nil
}

// Load loads the configuration from a base config file defined by a
// path and format.
func (l ConfigLoader) Load(id string, path string, format string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error while parsing the list of sources")
		}
	}()

	source, err := l.sourceFactory.Create(ConfigSourceTypeFile, path, format)
	if err != nil {
		return err
	}
	if err = l.config.AddSource(id, 0, source); err != nil {
		return err
	}

	entries := l.config.Get("config.sources")
	if entries == nil {
		return nil
	}

	for _, conf := range entries.([]interface{}) {
		if err = l.loadSource(conf.(ConfigPartial)); err != nil {
			return err
		}
	}

	return nil
}

func (l ConfigLoader) loadSource(conf ConfigPartial) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	id := conf.String("id")
	priority := conf.Int("priority")

	var source ConfigSource
	if source, err = l.sourceFactory.CreateConfig(conf); err != nil {
		return err
	}

	if err = l.config.AddSource(id, priority, source); err != nil {
		return err
	}

	return nil
}

/// ---------------------------------------------------------------------------
/// ConfigParams
/// ---------------------------------------------------------------------------

// ConfigParams defines the config provider parameters storing structure
// that will be needed when instantiating a new provider
type ConfigParams struct {
	ConfigID         string
	FileSystemID     string
	SourceFactoryID  string
	DecoderFactoryID string
	LoaderID         string
	ObserveFrequency time.Duration
	BaseSourceActive bool
	BaseSourceID     string
	BaseSourcePath   string
	BaseSourceFormat string
}

// NewConfigParams creates a new config provider
// parameters instance with the default values.
func NewConfigParams() *ConfigParams {
	configID := ContainerConfigID
	if env := os.Getenv(EnvContainerConfigID); env != "" {
		configID = env
	}

	fileSystemID := ContainerFileSystemID
	if env := os.Getenv(EnvContainerFileSystemID); env != "" {
		fileSystemID = env
	}

	sourceFactoryID := ContainerConfigSourceFactoryID
	if env := os.Getenv(EnvContainerConfigSourceFactoryID); env != "" {
		sourceFactoryID = env
	}

	decoderFactoryID := ContainerConfigDecoderFactoryID
	if env := os.Getenv(EnvContainerConfigDecoderFactoryID); env != "" {
		decoderFactoryID = env
	}

	loaderID := ContainerConfigLoaderID
	if env := os.Getenv(EnvContainerConfigLoaderID); env != "" {
		loaderID = env
	}

	observeFrequency := ConfigObserveFrequency
	if env := os.Getenv(EnvConfigObserveFrequency); env != "" {
		seconds, _ := strconv.Atoi(env)
		observeFrequency = time.Second * time.Duration(seconds)
	}

	baseSourceActive := ConfigBaseSourceActive
	if env := os.Getenv(EnvConfigBaseSourceActive); env != "" {
		baseSourceActive = env == "true"
	}

	baseSourceID := ConfigBaseSourceID
	if env := os.Getenv(EnvConfigBaseSourceID); env != "" {
		baseSourceID = env
	}

	baseSourcePath := ConfigBaseSourcePath
	if env := os.Getenv(EnvConfigBaseSourcePath); env != "" {
		baseSourcePath = env
	}

	baseSourceFormat := ConfigBaseSourceFormat
	if env := os.Getenv(EnvConfigBaseSourceFormat); env != "" {
		baseSourceFormat = env
	}

	return &ConfigParams{
		ConfigID:         configID,
		FileSystemID:     fileSystemID,
		SourceFactoryID:  sourceFactoryID,
		DecoderFactoryID: decoderFactoryID,
		LoaderID:         loaderID,
		ObserveFrequency: observeFrequency,
		BaseSourceActive: baseSourceActive,
		BaseSourceID:     baseSourceID,
		BaseSourcePath:   baseSourcePath,
		BaseSourceFormat: baseSourceFormat,
	}
}

/// ---------------------------------------------------------------------------
/// ConfigProvider
/// ---------------------------------------------------------------------------

// ConfigProvider defines the default configuration provider to be used on
// the application initialization to register the configuration services.
type ConfigProvider struct {
	params *ConfigParams
}

// NewConfigProvider will create a new configuration provider instance used to
// register the basic configuration objects in the application container.
func NewConfigProvider(params *ConfigParams) *ConfigProvider {
	if params == nil {
		params = NewConfigParams()
	}

	return &ConfigProvider{
		params: params,
	}
}

// Register will register the configuration section instances in the
// application container.
func (p ConfigProvider) Register(container *AppContainer) error {
	if container == nil {
		return fmt.Errorf("invalid nil 'container' argument")
	}

	_ = container.Add(p.params.DecoderFactoryID, func(container *AppContainer) (interface{}, error) {
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())

		return decoderFactory, nil
	})

	_ = container.Add(p.params.SourceFactoryID, func(container *AppContainer) (obj interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()

		fileSystem, err := container.Get(p.params.FileSystemID)
		if err != nil {
			return nil, err
		}

		decoderFactory, err := container.Get(p.params.DecoderFactoryID)
		if err != nil {
			return nil, err
		}

		fileSourceFactoryStrategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem.(afero.Fs), decoderFactory.(*ConfigDecoderFactory))
		observableFileSourceFactoryStrategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem.(afero.Fs), decoderFactory.(*ConfigDecoderFactory))

		sourceFactory := NewConfigSourceFactory()
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		return sourceFactory, nil
	})

	_ = container.Add(p.params.ConfigID, func(container *AppContainer) (interface{}, error) {
		return NewConfig(p.params.ObserveFrequency)
	})

	_ = container.Add(p.params.LoaderID, func(container *AppContainer) (obj interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()

		config, err := container.Get(p.params.ConfigID)
		if err != nil {
			return nil, err
		}

		sourceFactory, err := container.Get(p.params.SourceFactoryID)
		if err != nil {
			return nil, err
		}

		return NewConfigLoader(config.(*Config), sourceFactory.(*ConfigSourceFactory))
	})

	return nil
}

// Boot will start the configuration config instance by calling the
// configuration loader with the defined provider base entry information.
func (p ConfigProvider) Boot(container *AppContainer) (err error) {
	if p.params.BaseSourceActive {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()

		loader, err := container.Get(p.params.LoaderID)
		if err != nil {
			return err
		}

		return loader.(*ConfigLoader).Load(p.params.BaseSourceID, p.params.BaseSourcePath, p.params.BaseSourceFormat)
	}
	return nil
}
