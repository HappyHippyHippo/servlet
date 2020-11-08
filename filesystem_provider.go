package servlet

import (
	"github.com/spf13/afero"
)

// FileSystemProvider defines the default configuration provider to be used on
// the application initialization to register the file system adapter service.
type FileSystemProvider struct {
	params *FileSystemProviderParams
}

// NewFileSystemProvider instantiate a new system library provider that will register a
// file system adaptor in the application container.
func NewFileSystemProvider(params *FileSystemProviderParams) *FileSystemProvider {
	if params == nil {
		params = NewFileSystemProviderParams()
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
