package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_NewConfigProvider(t *testing.T) {
	t.Run("without params", func(t *testing.T) {
		if provider := NewConfigProvider(nil); provider == nil {
			t.Error("didn't returned a valid reference")
		} else if !reflect.DeepEqual(NewConfigProviderParams(), provider.params) {
			t.Errorf("stored the (%v) parameters", provider.params)
		}
	})

	t.Run("with defined params", func(t *testing.T) {
		params := NewConfigProviderParams()
		if provider := NewConfigProvider(params); provider == nil {
			t.Error("didn't returned a valid reference")
		} else if params != provider.params {
			t.Errorf("stored the (%v) parameters", provider.params)
		}
	})
}

func Test_ConfigProvider_Register(t *testing.T) {
	t.Run("nil container", func(t *testing.T) {
		provider := NewConfigProvider(nil)
		if err := provider.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'container' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register components", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		provider := NewConfigProvider(nil)

		if err := provider.Register(container); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !container.Has(ContainerConfigDecoderFactoryStrategyYamlID) {
			t.Errorf("didn't registered the config decoder factory strategy yaml : %v", provider)
		} else if !container.Has(ContainerConfigDecoderFactoryID) {
			t.Errorf("didn't registered the config decoder factory : %v", provider)
		} else if !container.Has(ContainerConfigSourceFactoryStrategyFileID) {
			t.Errorf("didn't registered the config source factory strategy file : %v", provider)
		} else if !container.Has(ContainerConfigSourceFactoryStrategyObservableFileID) {
			t.Errorf("didn't registered the config source factory strategy observable file : %v", provider)
		} else if !container.Has(ContainerConfigSourceFactoryStrategyEnvironmentID) {
			t.Errorf("didn't registered the config source factory strategy environment : %v", provider)
		} else if !container.Has(ContainerConfigSourceFactoryID) {
			t.Errorf("didn't registered the config source factory : %v", provider)
		} else if !container.Has(ContainerConfigID) {
			t.Errorf("didn't registered the config : %v", provider)
		} else if !container.Has(ContainerConfigLoaderID) {
			t.Errorf("didn't registered the config loader : %v", provider)
		}
	})

	t.Run("retrieving config yaml decoder factory strategy", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewConfigProvider(nil).Register(container)

		if strategy, err := container.Get(ContainerConfigDecoderFactoryStrategyYamlID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if strategy == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch strategy.(type) {
			case *ConfigDecoderFactoryStrategyYaml:
			default:
				t.Error("didn't returned a yaml decoder factory strategy reference")
			}
		}
	})

	t.Run("retrieving config decoder factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewConfigProvider(nil).Register(container)

		if factory, err := container.Get(ContainerConfigDecoderFactoryID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch factory.(type) {
			case *ConfigDecoderFactory:
			default:
				t.Error("didn't returned a decoder factory reference")
			}
		}
	})

	t.Run("error retrieving file system on retrieving the source factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if strategy, err := container.Get(ContainerConfigSourceFactoryStrategyFileID); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid file system on retrieving the source factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if strategy, err := container.Get(ContainerConfigSourceFactoryStrategyFileID); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving decoder factory on retrieving the source factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigDecoderFactoryID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if strategy, err := container.Get(ContainerConfigSourceFactoryStrategyFileID); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid decoder factory on retrieving the source factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigDecoderFactoryID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if strategy, err := container.Get(ContainerConfigSourceFactoryStrategyFileID); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("retrieving the source factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		if strategy, err := container.Get(ContainerConfigSourceFactoryStrategyFileID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if strategy == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch strategy.(type) {
			case *ConfigSourceFactoryStrategyFile:
			default:
				t.Error("didn't returned a source factory strategy file reference")
			}
		}
	})

	t.Run("error retrieving file system on retrieving the source factory strategy observable file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if strategy, err := container.Get(ContainerConfigSourceFactoryStrategyObservableFileID); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid file system on retrieving the source factory strategy observable file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if strategy, err := container.Get(ContainerConfigSourceFactoryStrategyObservableFileID); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving decoder factory on retrieving the source factory strategy observable file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigDecoderFactoryID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if strategy, err := container.Get(ContainerConfigSourceFactoryStrategyObservableFileID); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid decoder factory on retrieving the source factory strategy observable file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigDecoderFactoryID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if strategy, err := container.Get(ContainerConfigSourceFactoryStrategyObservableFileID); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("retrieving the source factory strategy observable file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		if strategy, err := container.Get(ContainerConfigSourceFactoryStrategyObservableFileID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if strategy == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch strategy.(type) {
			case *ConfigSourceFactoryStrategyObservableFile:
			default:
				t.Error("didn't returned a source factory strategy observable file reference")
			}
		}
	})

	t.Run("retrieving the source factory strategy environment", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewConfigProvider(nil).Register(container)

		if strategy, err := container.Get(ContainerConfigSourceFactoryStrategyEnvironmentID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if strategy == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch strategy.(type) {
			case *ConfigSourceFactoryStrategyEnvironment:
			default:
				t.Error("didn't returned a source factory strategy environment reference")
			}
		}
	})

	t.Run("retrieving config source factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewConfigProvider(nil).Register(container)

		if factory, err := container.Get(ContainerConfigSourceFactoryID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch factory.(type) {
			case *ConfigSourceFactory:
			default:
				t.Error("didn't returned a source factory reference")
			}
		}
	})

	t.Run("retrieving config", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		if config, err := container.Get(ContainerConfigID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if config == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch config.(type) {
			case *Config:
			default:
				t.Error("didn't returned a config reference")
			}
		}
	})

	t.Run("error retrieving config on retrieving loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if loader, err := container.Get(ContainerConfigLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid config on retrieving loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if loader, err := container.Get(ContainerConfigLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving config source factory on retrieving loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigSourceFactoryID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if loader, err := container.Get(ContainerConfigLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid config source factory on retrieving loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigSourceFactoryID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if loader, err := container.Get(ContainerConfigLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("retrieving config loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		if loader, err := container.Get(ContainerConfigLoaderID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if loader == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch loader.(type) {
			case *ConfigLoader:
			default:
				t.Error("didn't returned a loader reference")
			}
		}
	})
}

func Test_ConfigProvider_Boot(t *testing.T) {
	t.Run("error retrieving config decoder factory", func(t *testing.T) {
		expected := fmt.Errorf("error")

		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigDecoderFactoryID, func(container *AppContainer) (interface{}, error) {
			return nil, expected
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if err != expected {
			t.Errorf("returned the unexpected (%v) error", err)
		}
	})

	t.Run("retrieving invalid config decoder factory", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigDecoderFactoryID, func(container *AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving config decoder factory strategy yaml", func(t *testing.T) {
		expected := fmt.Errorf("error")

		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigDecoderFactoryStrategyYamlID, func(container *AppContainer) (interface{}, error) {
			return nil, expected
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if err != expected {
			t.Errorf("returned the unexpected (%v) error", err)
		}
	})

	t.Run("retrieving invalid config decoder factory strategy yaml", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigDecoderFactoryStrategyYamlID, func(container *AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving config source factory", func(t *testing.T) {
		expected := fmt.Errorf("error")

		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigSourceFactoryID, func(container *AppContainer) (interface{}, error) {
			return nil, expected
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if err != expected {
			t.Errorf("returned the unexpected (%v) error", err)
		}
	})

	t.Run("retrieving invalid config source factory", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = NewFileSystemProvider(nil).Register(container)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigSourceFactoryID, func(container *AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving config source factory strategy file", func(t *testing.T) {
		expected := fmt.Errorf("error")

		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = NewFileSystemProvider(nil).Register(container)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigSourceFactoryStrategyFileID, func(container *AppContainer) (interface{}, error) {
			return nil, expected
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if err != expected {
			t.Errorf("returned the unexpected (%v) error", err)
		}
	})

	t.Run("retrieving invalid config source factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = NewFileSystemProvider(nil).Register(container)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigSourceFactoryStrategyFileID, func(container *AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving config source factory strategy observable file", func(t *testing.T) {
		expected := fmt.Errorf("error")

		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = NewFileSystemProvider(nil).Register(container)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigSourceFactoryStrategyObservableFileID, func(container *AppContainer) (interface{}, error) {
			return nil, expected
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if err != expected {
			t.Errorf("returned the unexpected (%v) error", err)
		}
	})

	t.Run("retrieving invalid config source factory strategy observable file", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = NewFileSystemProvider(nil).Register(container)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigSourceFactoryStrategyObservableFileID, func(container *AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving config source factory strategy environment", func(t *testing.T) {
		expected := fmt.Errorf("error")

		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = NewFileSystemProvider(nil).Register(container)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigSourceFactoryStrategyEnvironmentID, func(container *AppContainer) (interface{}, error) {
			return nil, expected
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if err != expected {
			t.Errorf("returned the unexpected (%v) error", err)
		}
	})

	t.Run("retrieving invalid config source factory strategy environment", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = NewFileSystemProvider(nil).Register(container)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigSourceFactoryStrategyEnvironmentID, func(container *AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("no entry source active", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)

		params := NewConfigProviderParams()
		params.EntrySourceActive = false
		provider := NewConfigProvider(params)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigLoaderID, func(container *AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if err := provider.Boot(container); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		provider := NewConfigProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigLoaderID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		provider := NewConfigProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigLoaderID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("add entry source into the config", func(t *testing.T) {
		content := "field: value"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(ConfigEntrySourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		container := NewAppContainer()
		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return fileSystem, nil
		})

		provider := NewConfigProvider(nil)
		_ = provider.Register(container)

		if err := provider.Boot(container); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}
