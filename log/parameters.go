package log

import (
	"os"

	"github.com/happyhippyhippo/servlet/config"
	"github.com/happyhippyhippo/servlet/sys"
)

const (
	// ContainerLoggerID defines the id to be used as the default of a
	// logger instance in the application container.
	ContainerLoggerID = "servlet.logger"

	// EnvContainerLoggerID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// logger id.
	EnvContainerLoggerID = "SERVLET_CONTAINER_LOGGER_ID"

	// ContainerFormatterFactoryID defines the id to be used as the default of
	// a logger formatter factory instance in the application container.
	ContainerFormatterFactoryID = "logger.factory.formatter"

	// EnvContainerFormatterFactoryID defines the name of the environment
	// variable to be checked for a overriding value for the application
	// container logger formatter factory id.
	EnvContainerFormatterFactoryID = "SERVLET_CONTAINER_LOGGER_FORMATTER_FACTORY_ID"

	// ContainerStreamFactoryID defines the id to be used as the default of a
	// logger source factory instance in the application container.
	ContainerStreamFactoryID = "logger.factory.stream"

	// EnvContainerStreamFactoryID defines the name of the environment
	// variable to be checked for a overriding value for the application
	// container logger stream factory id.
	EnvContainerStreamFactoryID = "SERVLET_CONTAINER_LOGGER_STREAM_FACTORY_ID"

	// ContainerLoaderID defines the id to be used as the default of a
	// logger loader instance in the application container.
	ContainerLoaderID = "servlet.logger.loader"

	// EnvContainerLoaderID defines the name of the environment
	// variable to be checked for a overriding value for the application
	// container logger loader id.
	EnvContainerLoaderID = "SERVLET_CONTAINER_LOGGER_LOADER_ID"
)

// Parameters defines the logging provider parameters storing structure
// that will be needed when instantiating a new provider
type Parameters struct {
	LoggerID           string
	FileSystemID       string
	ConfigID           string
	FormatterFactoryID string
	StreamFactoryID    string
	LoaderID           string
}

// NewParameters will instantiate a new log provider parameters
// storing instance.
func NewParameters(
	loggerID string,
	fileSystemID string,
	configID string,
	formatterFactoryID string,
	streamFactoryID string,
	loaderID string) Parameters {

	return Parameters{
		LoggerID:           loggerID,
		FileSystemID:       fileSystemID,
		ConfigID:           configID,
		FormatterFactoryID: formatterFactoryID,
		StreamFactoryID:    streamFactoryID,
		LoaderID:           loaderID,
	}
}

// NewDefaultParameters will instantiate a new log provider parameters
// storing instance with the servlet default values.
func NewDefaultParameters() Parameters {
	loggerID := ContainerLoggerID
	if env := os.Getenv(EnvContainerLoggerID); env != "" {
		loggerID = env
	}

	fileSystemID := sys.ContainerFileSystemID
	if env := os.Getenv(sys.EnvContainerFileSystemID); env != "" {
		fileSystemID = env
	}

	configID := config.ContainerConfigID
	if env := os.Getenv(config.EnvContainerConfigID); env != "" {
		configID = env
	}

	formatterFactoryID := ContainerFormatterFactoryID
	if env := os.Getenv(EnvContainerFormatterFactoryID); env != "" {
		formatterFactoryID = env
	}

	streamFactoryID := ContainerStreamFactoryID
	if env := os.Getenv(EnvContainerStreamFactoryID); env != "" {
		streamFactoryID = env
	}

	loaderID := ContainerLoaderID
	if env := os.Getenv(EnvContainerLoaderID); env != "" {
		loaderID = env
	}

	return Parameters{
		LoggerID:           loggerID,
		FileSystemID:       fileSystemID,
		ConfigID:           configID,
		FormatterFactoryID: formatterFactoryID,
		StreamFactoryID:    streamFactoryID,
		LoaderID:           loaderID,
	}
}
