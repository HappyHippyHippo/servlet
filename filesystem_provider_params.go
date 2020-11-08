package servlet

import "os"

// FileSystemProviderParams defines the system provider parameters storing structure
// that will be needed when instantiating a new provider
type FileSystemProviderParams struct {
	FileSystemID string
}

// NewFileSystemProviderParams instantiate a new file system
// provider parameters object with the default values.
func NewFileSystemProviderParams() *FileSystemProviderParams {
	params := &FileSystemProviderParams{
		FileSystemID: ContainerFileSystemID,
	}

	if env := os.Getenv(EnvContainerFileSystemID); env != "" {
		params.FileSystemID = env
	}

	return params
}
