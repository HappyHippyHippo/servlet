package servlet

import (
	"github.com/spf13/afero"
	"reflect"
	"testing"
)

func Test_NewFileSystemProvider(t *testing.T) {
	t.Run("without params", func(t *testing.T) {
		if provider := NewFileSystemProvider(nil); provider == nil {
			t.Error("didn't returned a valid reference")
		} else if !reflect.DeepEqual(NewFileSystemProviderParams(), provider.params) {
			t.Errorf("stored the (%v) parameters", provider.params)
		}
	})

	t.Run("with defined params", func(t *testing.T) {
		params := NewFileSystemProviderParams()
		if provider := NewFileSystemProvider(params); provider == nil {
			t.Error("didn't returned a valid reference")
		} else if params != provider.params {
			t.Errorf("stored the (%v) parameters", provider.params)
		}
	})
}

func Test_FileSystemProvider_Register(t *testing.T) {
	a := NewApp()

	p := NewFileSystemProvider(nil)
	_ = p.Register(a.container)

	t.Run("register the file system", func(t *testing.T) {
		if f, ok := a.container.factories[ContainerFileSystemID]; !ok {
			t.Error("didn't registered the file system in the application container")
		} else {
			e, _ := f(a.container)
			switch e.(type) {
			case *afero.OsFs:
			default:
				t.Error("didn't returned the file system form the container")
			}
		}
	})
}

func Test_FileSystemProvider_Boot(t *testing.T) {
	a := NewApp()

	p := NewFileSystemProvider(nil)
	_ = p.Register(a.container)

	if err := p.Boot(a.container); err != nil {
		t.Errorf("returned the (%v) error", err)
	}
}
