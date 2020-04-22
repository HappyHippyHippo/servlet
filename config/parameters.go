package config

import (
	"os"
	"strconv"
	"time"

	"github.com/happyhippyhippo/servlet/sys"
)

const (
	// ContainerConfigID defines the id to be used as the default of a
	// config instance in the application container.
	ContainerConfigID = "servlet.config"

	// EnvContainerConfigID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// config id.
	EnvContainerConfigID = "SERVLET_CONTAINER_CONFIG_ID"

	// ContainerDecoderFactoryID defines the id to be used as the default of a
	// config decoder factory instance in the application container.
	ContainerDecoderFactoryID = "servlet.config.factory.decoder"

	// EnvContainerDecoderFactoryID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// config decoder factory id.
	EnvContainerDecoderFactoryID = "SERVLET_CONTAINER_CONFIG_DECODER_FACTORY_ID"

	// ContainerSourceFactoryID defines the id to be used as the default of a
	// config source factory instance in the application container.
	ContainerSourceFactoryID = "servlet.config.factory.source"

	// EnvContainerSourceFactoryID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// config source factory id.
	EnvContainerSourceFactoryID = "SERVLET_CONTAINER_CONFIG_SOURCE_FACTORY_ID"

	// ContainerLoaderID defines the id to be used as the default of a
	// config loader instance in the application container.
	ContainerLoaderID = "servlet.config.loader"

	// EnvContainerLoaderID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// config loaded id.
	EnvContainerLoaderID = "SERVLET_CONTAINER_CONFIG_LOADER_ID"

	// ObserveFrequency defines the id to be used as the default of a
	// config observable source frequancy time.
	ObserveFrequency = time.Second * 0

	// EnvObserveFrequency defines the name of the environment variable
	// to be checked for a overriding value for the config observe frequency.
	EnvObserveFrequency = "SERVLET_CONFIG_OBSERVE_FREQUENCY"

	// BaseSourceID defines the id to be used as the default of the
	// base config source id to be used as the loader entry.
	BaseSourceID = "base"

	// EnvBaseSourceID defines the name of the environment variable
	// to be checked for a overriding value for the config base source id.
	EnvBaseSourceID = "SERVLET_CONFIG_BASE_SOURCE_ID"

	// BaseSourcePath defines the base config source path
	// to be used as the loader entry.
	BaseSourcePath = "config/config.yaml"

	// EnvBaseSourcePath defines the name of the environment variable
	// to be checked for a overriding value for the config base source path.
	EnvBaseSourcePath = "SERVLET_CONFIG_BASE_SOURCE_PATH"

	// BaseSourceFormat defines the base config source format
	// to be used as the loader entry.
	BaseSourceFormat = DecoderFormatYAML

	// EnvBaseSourceFormat defines the name of the environment variable
	// to be checked for a overriding value for the config base source format.
	EnvBaseSourceFormat = "SERVLET_CONFIG_BASE_SOURCE_FORMAT"
)

// Parameters defines the config provider parameters storing structure
// that will be needed when instantiating a new provider
type Parameters struct {
	ConfigID         string
	FileSystemID     string
	SourceFactoryID  string
	DecoderFactoryID string
	LoaderID         string
	ObserveFrequency time.Duration
	BaseSourceID     string
	BaseSourcePath   string
	BaseSourceFormat string
}

// NewParameters creates a new config provider parameters instance.
func NewParameters(
	configID string,
	fileSystemID string,
	sourceFactoryID string,
	decoderFactoryID string,
	loaderID string,
	observeFrequency time.Duration,
	baseSourceID string,
	baseSourcePath string,
	baseSourceFormat string) Parameters {

	return Parameters{
		ConfigID:         configID,
		FileSystemID:     fileSystemID,
		SourceFactoryID:  sourceFactoryID,
		DecoderFactoryID: decoderFactoryID,
		LoaderID:         loaderID,
		ObserveFrequency: observeFrequency,
		BaseSourceID:     baseSourceID,
		BaseSourcePath:   baseSourcePath,
		BaseSourceFormat: baseSourceFormat,
	}
}

// NewDefaultParameters creates a new config provider
// parameters instance with the default values.
func NewDefaultParameters() Parameters {
	configID := ContainerConfigID
	if env := os.Getenv(EnvContainerConfigID); env != "" {
		configID = env
	}

	fileSystemID := sys.ContainerFileSystemID
	if env := os.Getenv(sys.EnvContainerFileSystemID); env != "" {
		fileSystemID = env
	}

	sourceFactoryID := ContainerSourceFactoryID
	if env := os.Getenv(EnvContainerSourceFactoryID); env != "" {
		sourceFactoryID = env
	}

	decoderFactoryID := ContainerDecoderFactoryID
	if env := os.Getenv(EnvContainerDecoderFactoryID); env != "" {
		decoderFactoryID = env
	}

	loaderID := ContainerLoaderID
	if env := os.Getenv(EnvContainerLoaderID); env != "" {
		loaderID = env
	}

	observeFrequency := ObserveFrequency
	if env := os.Getenv(EnvObserveFrequency); env != "" {
		seconds, _ := strconv.Atoi(env)
		observeFrequency = time.Second * time.Duration(seconds)
	}

	baseSourceID := BaseSourceID
	if env := os.Getenv(EnvBaseSourceID); env != "" {
		baseSourceID = env
	}

	baseSourcePath := BaseSourcePath
	if env := os.Getenv(EnvBaseSourcePath); env != "" {
		baseSourcePath = env
	}

	baseSourceFormat := BaseSourceFormat
	if env := os.Getenv(EnvBaseSourceFormat); env != "" {
		baseSourceFormat = env
	}

	return Parameters{
		ConfigID:         configID,
		FileSystemID:     fileSystemID,
		SourceFactoryID:  sourceFactoryID,
		DecoderFactoryID: decoderFactoryID,
		LoaderID:         loaderID,
		ObserveFrequency: observeFrequency,
		BaseSourceID:     baseSourceID,
		BaseSourcePath:   baseSourcePath,
		BaseSourceFormat: baseSourceFormat,
	}
}
