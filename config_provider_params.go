package servlet

import (
	"os"
	"strconv"
	"time"
)

// ConfigProviderParams defines the config provider parameters storing structure
// that will be needed when instantiating a new provider
type ConfigProviderParams struct {
	ConfigID                              string
	FileSystemID                          string
	SourceFactoryStrategyFileID           string
	SourceFactoryStrategyObservableFileID string
	SourceFactoryStrategyEnvironmentID    string
	SourceFactoryID                       string
	DecoderFactoryStrategyYamlID          string
	DecoderFactoryID                      string
	LoaderID                              string
	ObserveFrequency                      time.Duration
	EntrySourceActive                     bool
	EntrySourceID                         string
	EntrySourcePath                       string
	EntrySourceFormat                     string
}

// NewConfigProviderParams creates a new config provider
// parameters instance with the default values.
func NewConfigProviderParams() *ConfigProviderParams {
	params := &ConfigProviderParams{
		ConfigID:                              ContainerConfigID,
		FileSystemID:                          ContainerFileSystemID,
		SourceFactoryStrategyFileID:           ContainerConfigSourceFactoryStrategyFileID,
		SourceFactoryStrategyObservableFileID: ContainerConfigSourceFactoryStrategyObservableFileID,
		SourceFactoryStrategyEnvironmentID:    ContainerConfigSourceFactoryStrategyEnvironmentID,
		SourceFactoryID:                       ContainerConfigSourceFactoryID,
		DecoderFactoryStrategyYamlID:          ContainerConfigDecoderFactoryStrategyYamlID,
		DecoderFactoryID:                      ContainerConfigDecoderFactoryID,
		LoaderID:                              ContainerConfigLoaderID,
		ObserveFrequency:                      ConfigObserveFrequency,
		EntrySourceActive:                     ConfigEntrySourceActive,
		EntrySourceID:                         ConfigEntrySourceID,
		EntrySourcePath:                       ConfigEntrySourcePath,
		EntrySourceFormat:                     ConfigEntrySourceFormat,
	}

	if env := os.Getenv(EnvContainerConfigID); env != "" {
		params.ConfigID = env
	}

	if env := os.Getenv(EnvContainerFileSystemID); env != "" {
		params.FileSystemID = env
	}

	if env := os.Getenv(EnvContainerConfigSourceFactoryStrategyFileID); env != "" {
		params.SourceFactoryStrategyFileID = env
	}

	if env := os.Getenv(EnvContainerConfigSourceFactoryStrategyObservableFileID); env != "" {
		params.SourceFactoryStrategyObservableFileID = env
	}

	if env := os.Getenv(EnvContainerConfigSourceFactoryStrategyEnvironmentID); env != "" {
		params.SourceFactoryStrategyEnvironmentID = env
	}

	if env := os.Getenv(EnvContainerConfigSourceFactoryID); env != "" {
		params.SourceFactoryID = env
	}

	if env := os.Getenv(EnvContainerConfigDecoderFactoryStrategyYamlID); env != "" {
		params.DecoderFactoryStrategyYamlID = env
	}

	if env := os.Getenv(EnvContainerConfigDecoderFactoryID); env != "" {
		params.DecoderFactoryID = env
	}

	if env := os.Getenv(EnvContainerConfigLoaderID); env != "" {
		params.LoaderID = env
	}

	if env := os.Getenv(EnvConfigObserveFrequency); env != "" {
		seconds, _ := strconv.Atoi(env)
		params.ObserveFrequency = time.Second * time.Duration(seconds)
	}

	if env := os.Getenv(EnvConfigEntrySourceActive); env != "" {
		params.EntrySourceActive = env == "true"
	}

	if env := os.Getenv(EnvConfigEntrySourceID); env != "" {
		params.EntrySourceID = env
	}

	if env := os.Getenv(EnvConfigEntrySourcePath); env != "" {
		params.EntrySourcePath = env
	}

	if env := os.Getenv(EnvConfigEntrySourceFormat); env != "" {
		params.EntrySourceFormat = env
	}

	return params
}
