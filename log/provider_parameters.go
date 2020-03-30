package log

import (
	"github.com/happyhippyhippo/servlet/config"
	"github.com/happyhippyhippo/servlet/sys"
)

const (
	// ContainerID defines the id to be used as the default of a
	// logger instance in the application container.
	ContainerID = "logger"

	// ContainerLoaderID defines the id to be used as the default of a
	// logger loader instance in the application container.
	ContainerLoaderID = "logger.loader"

	// ContainerFormatterFactoryID defines the id to be used as the default of
	// a logger formatter factory instance in the application container.
	ContainerFormatterFactoryID = "logger.factory.formatter"

	// ContainerStreamFactoryID defines the id to be used as the default of a
	// logger source factory instance in the application container.
	ContainerStreamFactoryID = "logger.factory.stream"
)

// ProviderParameters interface defines the methods of a log provider
// parameters storing instance.
type ProviderParameters interface {
	GetID() string
	SetID(id string) ProviderParameters
	GetFileSystemID() string
	SetFileSystemID(id string) ProviderParameters
	GetConfigID() string
	SetConfigID(id string) ProviderParameters
	GetFormatterFactoryID() string
	SetFormatterFactoryID(id string) ProviderParameters
	GetStreamFactoryID() string
	SetStreamFactoryID(id string) ProviderParameters
	GetLoaderID() string
	SetLoaderID(id string) ProviderParameters
}

type providerParameters struct {
	id                 string
	fileSystemID       string
	configID           string
	formatterFactoryID string
	streamFactoryID    string
	loaderID           string
}

// NewProviderParameters will instantiate a new log provider parameters
// storing instance.
func NewProviderParameters(
	id string,
	fileSystemID string,
	configID string,
	formatterFactoryID string,
	streamFactoryID string,
	loaderID string) ProviderParameters {

	return &providerParameters{
		id:                 id,
		fileSystemID:       fileSystemID,
		configID:           configID,
		formatterFactoryID: formatterFactoryID,
		streamFactoryID:    streamFactoryID,
		loaderID:           loaderID,
	}
}

// NewDefaultProviderParameters will instantiate a new log provider parameters
// storing instance with the servlet default values.
func NewDefaultProviderParameters() ProviderParameters {
	return &providerParameters{
		id:                 ContainerID,
		fileSystemID:       sys.ContainerFileSystemID,
		configID:           config.ContainerID,
		formatterFactoryID: ContainerFormatterFactoryID,
		streamFactoryID:    ContainerStreamFactoryID,
		loaderID:           ContainerLoaderID,
	}
}

// GetID retrieves the id used to register the logger instance in the
// application container.
func (p providerParameters) GetID() string {
	return p.id
}

// SetID updates the id used to register the logger instance in the
// application container.
func (p *providerParameters) SetID(id string) ProviderParameters {
	p.id = id
	return p
}

// GetFileSystemID retrieves the id used to register the
// file system adapter instance in the application container.
func (p providerParameters) GetFileSystemID() string {
	return p.fileSystemID
}

// SetFileSystemID updates the id used to register the
// file system adapter instance in the application container.
func (p *providerParameters) SetFileSystemID(id string) ProviderParameters {
	p.fileSystemID = id
	return p
}

// GetConfigID retrieves the id used to register the
// config instance in the application container.
func (p providerParameters) GetConfigID() string {
	return p.configID
}

// SetConfigID updates the id used to register the
// config instance in the application container.
func (p *providerParameters) SetConfigID(id string) ProviderParameters {
	p.configID = id
	return p
}

// GetFormatterFactoryID retrieves the id used to register the
// logging formatter factory instance in the application container.
func (p providerParameters) GetFormatterFactoryID() string {
	return p.formatterFactoryID
}

// SetFormatterFactoryID updates the id used to register the
// logging formatter factory instance in the application container.
func (p *providerParameters) SetFormatterFactoryID(id string) ProviderParameters {
	p.formatterFactoryID = id
	return p
}

// GetStreamFactoryID retrieves the id used to register the
// logging stream factory instance in the application container.
func (p providerParameters) GetStreamFactoryID() string {
	return p.streamFactoryID
}

// SetStreamFactoryID updates the id used to register the
// logging stream factory instance in the application container.
func (p *providerParameters) SetStreamFactoryID(id string) ProviderParameters {
	p.streamFactoryID = id
	return p
}

// GetLoaderID retrieves the id used to register the
// log loader instance in the application container.
func (p providerParameters) GetLoaderID() string {
	return p.loaderID
}

// SetLoaderID updates the id used to register the
// log loader instance in the application container.
func (p *providerParameters) SetLoaderID(id string) ProviderParameters {
	p.loaderID = id
	return p
}
