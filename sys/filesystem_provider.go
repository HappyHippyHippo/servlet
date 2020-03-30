package sys

import (
	"github.com/happyhippyhippo/servlet"
	"github.com/spf13/afero"
)

type fileSystemProvider struct {
	id string
}

// NewFileSystemProvider instantiate a new system provider.
func NewFileSystemProvider(p FileSystemProviderParameters) servlet.Provider {
	if p == nil {
		p = NewDefaultFileSystemProviderParameters()
	}

	return &fileSystemProvider{
		id: p.GetID(),
	}
}

// Register will register in the container a new file system adapter instance.
func (p fileSystemProvider) Register(c servlet.Container) {
	c.Add(p.id, func(c servlet.Container) interface{} {
		return afero.NewOsFs()
	})
}

// Boot (no-op).
func (fileSystemProvider) Boot(c servlet.Container) {
}
