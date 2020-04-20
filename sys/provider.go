package sys

import (
	"github.com/happyhippyhippo/servlet"
	"github.com/spf13/afero"
)

type provider struct {
	params Parameters
}

// NewProvider instantiate a new system library provider that will register a
// file system adaptor in the application container.
func NewProvider(params Parameters) servlet.Provider {
	return &provider{
		params: params,
	}
}

// Register will add to the container a new file system adapter instance.
func (p provider) Register(container servlet.Container) {
	container.Add(p.params.FileSystemID, func(container servlet.Container) interface{} {
		return afero.NewOsFs()
	})
}

// Boot (no-op).
func (provider) Boot(c servlet.Container) {
}
