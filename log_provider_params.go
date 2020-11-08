package servlet

import "os"

// LogProviderParams defines the logging provider parameters storing structure
// that will be needed when instantiating a new provider
type LogProviderParams struct {
	LoggerID                       string
	FileSystemID                   string
	ConfigID                       string
	FormatterFactoryStrategyJSONID string
	FormatterFactoryID             string
	StreamFactoryStrategyFileID    string
	StreamFactoryID                string
	LoaderID                       string
}

// NewLogProviderParams will instantiate a new log provider parameters
// storing instance with the servlet default values.
func NewLogProviderParams() *LogProviderParams {
	params := &LogProviderParams{
		LoggerID:                       ContainerLoggerID,
		FileSystemID:                   ContainerFileSystemID,
		ConfigID:                       ContainerConfigID,
		FormatterFactoryStrategyJSONID: ContainerLogFormatterFactoryStrategyJSONID,
		FormatterFactoryID:             ContainerLogFormatterFactoryID,
		StreamFactoryStrategyFileID:    ContainerLogStreamFactoryStrategyFileID,
		StreamFactoryID:                ContainerLogStreamFactoryID,
		LoaderID:                       ContainerLogLoaderID,
	}

	if env := os.Getenv(EnvContainerLoggerID); env != "" {
		params.LoggerID = env
	}

	if env := os.Getenv(EnvContainerFileSystemID); env != "" {
		params.FileSystemID = env
	}

	if env := os.Getenv(EnvContainerConfigID); env != "" {
		params.ConfigID = env
	}

	if env := os.Getenv(EnvContainerLogFormatterFactoryStrategyJSONID); env != "" {
		params.FormatterFactoryStrategyJSONID = env
	}

	if env := os.Getenv(EnvContainerLogFormatterFactoryID); env != "" {
		params.FormatterFactoryID = env
	}

	if env := os.Getenv(EnvContainerLogStreamFactoryStrategyFileID); env != "" {
		params.StreamFactoryStrategyFileID = env
	}

	if env := os.Getenv(EnvContainerLogStreamFactoryID); env != "" {
		params.StreamFactoryID = env
	}

	if env := os.Getenv(EnvContainerLogLoaderID); env != "" {
		params.LoaderID = env
	}

	return params
}
