package config

import (
	"os"
	"time"

	"github.com/happyhippyhippo/servlet/sys"
)

const (
	// ContainerID defines the id to be used as the default of a
	// config instance in the application container.
	ContainerID = "config"

	// ContainerLoaderID defines the id to be used as the default of a
	// config loader instance in the application container.
	ContainerLoaderID = "config.loader"

	// ContainerDecoderFactoryID defines the id to be used as the default of a
	// config decoder factory instance in the application container.
	ContainerDecoderFactoryID = "config.factory.decoder"

	// ContainerSourceFactoryID defines the id to be used as the default of a
	// config source factory instance in the application container.
	ContainerSourceFactoryID = "config.factory.source"

	// ContainerObserveFrequency defines the id to be used as the default of a
	// config observable source frequancy time.
	ContainerObserveFrequency = time.Second * 0

	// ContainerBaseSourceID defines the id to be used as the default of the
	// base config source id to be used as the loader entry.
	ContainerBaseSourceID = "config.sources.base"

	// ContainerBaseSourcePath defines the base config source path
	// to be used as the loader entry.
	ContainerBaseSourcePath = "config/config.yaml"

	// ContainerBaseSourceFormat defines the base config source format
	// to be used as the loader entry.
	ContainerBaseSourceFormat = DecoderFormatYAML

	// EnvironmentBaseSourcePath defines the name of the environment variable
	// to be checked for a overriding value for the base source path.
	EnvironmentBaseSourcePath = "SERVLET_BASE_SOURCE_PATH"

	// EnvironmentBaseSourceFormat defines the name of the environment variable
	// to be checked for a overriding value for the base source format.
	EnvironmentBaseSourceFormat = "SERVLET_BASE_SOURCE_FORMAT"
)

// ProviderParameters interface defines the methods of a config provider
// parameters storing instance.
type ProviderParameters interface {
	GetID() string
	SetID(id string) ProviderParameters
	GetFileSystemID() string
	SetFileSystemID(id string) ProviderParameters
	GetLoaderID() string
	SetLoaderID(id string) ProviderParameters
	GetSourceFactoryID() string
	SetSourceFactoryID(id string) ProviderParameters
	GetDecoderFactoryID() string
	SetDecoderFactoryID(id string) ProviderParameters
	GetObserveFrequency() time.Duration
	SetObserveFrequency(frequency time.Duration) ProviderParameters
	GetBaseSourceID() string
	SetBaseSourceID(id string) ProviderParameters
	GetBaseSourcePath() string
	SetBaseSourcePath(path string) ProviderParameters
	GetBaseSourceFormat() string
	SetBaseSourceFormat(format string) ProviderParameters
}

type providerParameters struct {
	id               string
	fileSystemID     string
	loaderID         string
	sourceFactoryID  string
	decoderFactoryID string
	observeFrequency time.Duration
	baseSourceID     string
	baseSourcePath   string
	baseSourceFormat string
}

// NewProviderParameters creates a new config provider parameters instance.
func NewProviderParameters(
	id string,
	fileSystemID string,
	loaderID string,
	sourceFactoryID string,
	decoderFactoryID string,
	observeFrequency time.Duration,
	baseSourceID string,
	baseSourcePath string,
	baseSourceFormat string) ProviderParameters {

	return &providerParameters{
		id:               id,
		fileSystemID:     fileSystemID,
		loaderID:         loaderID,
		sourceFactoryID:  sourceFactoryID,
		decoderFactoryID: decoderFactoryID,
		observeFrequency: observeFrequency,
		baseSourceID:     baseSourceID,
		baseSourcePath:   baseSourcePath,
		baseSourceFormat: baseSourceFormat,
	}
}

// NewDefaultProviderParameters creates a new config provider
// parameters instance with the default values.
func NewDefaultProviderParameters() ProviderParameters {
	baseSourcePath := ContainerBaseSourcePath
	if env := os.Getenv(EnvironmentBaseSourcePath); env != "" {
		baseSourcePath = env
	}

	baseSourceFormat := ContainerBaseSourceFormat
	if env := os.Getenv(EnvironmentBaseSourceFormat); env != "" {
		baseSourceFormat = env
	}

	return &providerParameters{
		id:               ContainerID,
		fileSystemID:     sys.ContainerFileSystemID,
		loaderID:         ContainerLoaderID,
		sourceFactoryID:  ContainerSourceFactoryID,
		decoderFactoryID: ContainerDecoderFactoryID,
		observeFrequency: ContainerObserveFrequency,
		baseSourceID:     ContainerBaseSourceID,
		baseSourcePath:   baseSourcePath,
		baseSourceFormat: baseSourceFormat,
	}
}

// GetID will retrieves the stored config instance container id.
func (p providerParameters) GetID() string {
	return p.id
}

// SetID will update the stored config instance container id.
func (p *providerParameters) SetID(id string) ProviderParameters {
	p.id = id
	return p
}

// GetFileSystemID will retrieves the stored file systen instance container id.
func (p providerParameters) GetFileSystemID() string {
	return p.fileSystemID
}

// SetFileSystemID will update the stored file systen instance container id.
func (p *providerParameters) SetFileSystemID(id string) ProviderParameters {
	p.fileSystemID = id
	return p
}

// GetLoaderID will retrieves the stored config loader
// instance container id.
func (p providerParameters) GetLoaderID() string {
	return p.loaderID
}

// SetLoaderID will update the stored config loader
// instance container id.
func (p *providerParameters) SetLoaderID(id string) ProviderParameters {
	p.loaderID = id
	return p
}

// GetSourceFactoryID will retrieves the stored config source
// factory instance container id.
func (p providerParameters) GetSourceFactoryID() string {
	return p.sourceFactoryID
}

// SetSourceFactoryID will update the stored config source
// factory instance container id.
func (p *providerParameters) SetSourceFactoryID(id string) ProviderParameters {
	p.sourceFactoryID = id
	return p
}

// GetDecoderFactoryID will retrieves the stored config decoder
// factory instance container id.
func (p providerParameters) GetDecoderFactoryID() string {
	return p.decoderFactoryID
}

// SetDecoderFactoryID will update the stored config decoder
// factory instance container id.
func (p *providerParameters) SetDecoderFactoryID(id string) ProviderParameters {
	p.decoderFactoryID = id
	return p
}

// GetObserveFrequency will retrieves the stored config
// source observing frequency.
func (p providerParameters) GetObserveFrequency() time.Duration {
	return p.observeFrequency
}

// SetObserveFrequency will update the stored config
// source observing frequency.
func (p *providerParameters) SetObserveFrequency(frequency time.Duration) ProviderParameters {
	p.observeFrequency = frequency
	return p
}

// GetBaseSourceID will retrieves the stored id to be used to register
// the config loader base entry in the config instance.
func (p providerParameters) GetBaseSourceID() string {
	return p.baseSourceID
}

// SetBaseSourceID will update the stored id to be used to register
// the config loader base entry in the config instance.
func (p *providerParameters) SetBaseSourceID(id string) ProviderParameters {
	p.baseSourceID = id
	return p
}

// GetBaseSourcePath will retrieves the stored path to be used to load
// the config base entry.
func (p providerParameters) GetBaseSourcePath() string {
	return p.baseSourcePath
}

// SetBaseSourcePath will update the stored path to be used to load
// the config base entry.
func (p *providerParameters) SetBaseSourcePath(path string) ProviderParameters {
	p.baseSourcePath = path
	return p
}

// GetBaseSourceFormat will retrieves the stored path to be used to load
// the config base entry.
func (p providerParameters) GetBaseSourceFormat() string {
	return p.baseSourceFormat
}

// SetBaseSourceFormat will update the stored path to be used to load
// the config base entry.
func (p *providerParameters) SetBaseSourceFormat(format string) ProviderParameters {
	p.baseSourceFormat = format
	return p
}
