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

	// ContainerObserveFrequency defines the id to be used as the default of a
	// config observable source frequancy time.
	ContainerObserveFrequency = time.Second * 0

	// EnvContainerObserveFrequency defines the name of the environment variable
	// to be checked for a overriding value for the config observe frequency.
	EnvContainerObserveFrequency = "SERVLET_CONTAINER_CONFIG_OBSERVE_FREQUENCY"

	// ContainerBaseSourceID defines the id to be used as the default of the
	// base config source id to be used as the loader entry.
	ContainerBaseSourceID = "servlet.config.sources.base"

	// EnvContainerBaseSourceID defines the name of the environment variable
	// to be checked for a overriding value for the config base source id.
	EnvContainerBaseSourceID = "SERVLET_CONTAINER_CONFIG_BASE_SOURCE_ID"

	// ContainerBaseSourcePath defines the base config source path
	// to be used as the loader entry.
	ContainerBaseSourcePath = "config/config.yaml"

	// EnvContainerBaseSourcePath defines the name of the environment variable
	// to be checked for a overriding value for the config base source path.
	EnvContainerBaseSourcePath = "SERVLET_CONTAINER_CONFIG_BASE_SOURCE_PATH"

	// ContainerBaseSourceFormat defines the base config source format
	// to be used as the loader entry.
	ContainerBaseSourceFormat = DecoderFormatYAML

	// EnvContainerBaseSourceFormat defines the name of the environment variable
	// to be checked for a overriding value for the config base source format.
	EnvContainerBaseSourceFormat = "SERVLET_CONTAINER_CONFIG_BASE_SOURCE_FORMAT"
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

	observeFrequency := ContainerObserveFrequency
	if env := os.Getenv(EnvContainerObserveFrequency); env != "" {
		seconds, _ := strconv.Atoi(env)
		observeFrequency = time.Second * time.Duration(seconds)
	}

	baseSourceID := ContainerBaseSourceID
	if env := os.Getenv(EnvContainerBaseSourceID); env != "" {
		baseSourceID = env
	}

	baseSourcePath := ContainerBaseSourcePath
	if env := os.Getenv(EnvContainerBaseSourcePath); env != "" {
		baseSourcePath = env
	}

	baseSourceFormat := ContainerBaseSourceFormat
	if env := os.Getenv(EnvContainerBaseSourceFormat); env != "" {
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
