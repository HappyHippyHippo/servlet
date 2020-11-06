package servlet

import (
	"github.com/spf13/afero"
)

import "os"

/// ---------------------------------------------------------------------------
/// constants
/// ---------------------------------------------------------------------------

const (
	// ContainerFileSystemID defines the default id used to register the
	// application file system adapter instance in the application container.
	ContainerFileSystemID = "servlet.filesystem"

	// EnvContainerFileSystemID defines the environment variable used to
	// override the default value for the container file system adapter id
	EnvContainerFileSystemID = "SERVLET_CONTAINER_FILE_SYSTEM_ID"
)

/// ---------------------------------------------------------------------------
/// FileSystemParams
/// ---------------------------------------------------------------------------

// FileSystemParams defines the system provider parameters storing structure
// that will be needed when instantiating a new provider
type FileSystemParams struct {
	FileSystemID string
}

// NewFileSystemParams instantiate a new file system
// provider parameters object with the default values.
func NewFileSystemParams() *FileSystemParams {
	fileSystemID := ContainerFileSystemID
	if env := os.Getenv(EnvContainerFileSystemID); env != "" {
		fileSystemID = env
	}

	return &FileSystemParams{
		FileSystemID: fileSystemID,
	}
}

/// ---------------------------------------------------------------------------
/// FileSystemProvider
/// ---------------------------------------------------------------------------

// FileSystemProvider defines the default configuration provider to be used on
// the application initialization to register the file system adapter service.
type FileSystemProvider struct {
	params *FileSystemParams
}

// NewFileSystemProvider instantiate a new system library provider that will register a
// file system adaptor in the application container.
func NewFileSystemProvider(params *FileSystemParams) *FileSystemProvider {
	if params == nil {
		params = NewFileSystemParams()
	}

	return &FileSystemProvider{
		params: params,
	}
}

// Register will add to the container a new file system adapter instance.
func (p FileSystemProvider) Register(c *AppContainer) error {
	return c.Add(p.params.FileSystemID, func(c *AppContainer) (interface{}, error) {
		return afero.NewOsFs(), nil
	})
}

// Boot (no-op).
func (FileSystemProvider) Boot(_ *AppContainer) error {
	return nil
}
