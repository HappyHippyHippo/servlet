package servlet

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_NewLogProvider(t *testing.T) {
	t.Run("without params", func(t *testing.T) {
		if provider := NewLogProvider(nil); provider == nil {
			t.Error("didn't returned a valid reference")
		} else if !reflect.DeepEqual(NewLogProviderParams(), provider.params) {
			t.Errorf("stored the (%v) parameters", provider.params)
		}
	})

	t.Run("with defined params", func(t *testing.T) {
		params := NewLogProviderParams()
		if provider := NewLogProvider(params); provider == nil {
			t.Error("didn't returned a valid reference")
		} else if params != provider.params {
			t.Errorf("stored the (%v) parameters", provider.params)
		}
	})
}

func Test_LogProvider_Register(t *testing.T) {
	t.Run("nil container", func(t *testing.T) {
		provider := NewLogProvider(nil)
		if err := provider.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'container' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register components", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewLogProvider(nil)

		if err := provider.Register(container); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !container.Has(ContainerLogFormatterFactoryStrategyJSONID) {
			t.Error("didn't registered the log formatter factory strategy json", err)
		} else if !container.Has(ContainerLogFormatterFactoryID) {
			t.Error("didn't registered the log formatter factory", err)
		} else if !container.Has(ContainerLogStreamFactoryStrategyFileID) {
			t.Error("didn't registered the log stream factory strategy file", err)
		} else if !container.Has(ContainerLogStreamFactoryID) {
			t.Error("didn't registered the log stream factory", err)
		} else if !container.Has(ContainerLoggerID) {
			t.Error("didn't registered the logger", err)
		} else if !container.Has(ContainerLogLoaderID) {
			t.Error("didn't registered the log loader", err)
		}
	})

	t.Run("retrieving log formatter factory strategy json", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewLogProvider(nil).Register(container)

		if strategy, err := container.Get(ContainerLogFormatterFactoryStrategyJSONID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if strategy == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch strategy.(type) {
			case *LogFormatterFactoryStrategyJSON:
			default:
				t.Error("didn't returned a formatter factory strategy json reference")
			}
		}
	})

	t.Run("retrieving log formatter factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewLogProvider(nil).Register(container)

		if factory, err := container.Get(ContainerLogFormatterFactoryID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch factory.(type) {
			case *LogFormatterFactory:
			default:
				t.Error("didn't returned a formatter factory reference")
			}
		}
	})

	t.Run("error retrieving file system on retrieving the stream factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if factory, err := container.Get(ContainerLogStreamFactoryStrategyFileID); factory != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid file system on retrieving the stream factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if factory, err := container.Get(ContainerLogStreamFactoryStrategyFileID); factory != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving formatter factory on retrieving the stream factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLogFormatterFactoryID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if factory, err := container.Get(ContainerLogStreamFactoryStrategyFileID); factory != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid formatter factory on retrieving the stream factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLogFormatterFactoryID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if factory, err := container.Get(ContainerLogStreamFactoryStrategyFileID); factory != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("retrieving log stream factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		if factory, err := container.Get(ContainerLogStreamFactoryID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch factory.(type) {
			case *LogStreamFactory:
			default:
				t.Error("didn't returned a stream factory reference")
			}
		}
	})

	t.Run("retrieving logger", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		if logger, err := container.Get(ContainerLoggerID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if logger == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch logger.(type) {
			case *Log:
			default:
				t.Error("didn't returned a logger reference")
			}
		}
	})

	t.Run("error retrieving logger on retrieving logger loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLoggerID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if loader, err := container.Get(ContainerLogLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid logger on retrieving logger loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLoggerID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if loader, err := container.Get(ContainerLogLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving stream factory on retrieving logger loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLogStreamFactoryID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if loader, err := container.Get(ContainerLogLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid source factory on retrieving logger loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLogStreamFactoryID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if loader, err := container.Get(ContainerLogLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("retrieving log loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		if loader, err := container.Get(ContainerLogLoaderID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if loader == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch loader.(type) {
			case *LogLoader:
			default:
				t.Error("didn't returned a loader reference")
			}
		}
	})
}

func Test_LogProvider_Boot(t *testing.T) {
	t.Run("error retrieving formatter factory", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerLogFormatterFactoryID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid formatter factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerLogFormatterFactoryID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving formatter factory strategy json", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerLogFormatterFactoryStrategyJSONID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid formatter factory strategy json", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerLogFormatterFactoryStrategyJSONID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving stream factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerLogStreamFactoryID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid stream factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerLogStreamFactoryID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving stream factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerLogStreamFactoryStrategyFileID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid stream factory strategy file", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerLogStreamFactoryStrategyFileID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerLogLoaderID, func(*AppContainer) (interface{}, error) {
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
		_ = NewConfigProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerLogLoaderID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving config", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		_ = container.Add(ContainerConfigID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("run boot log loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		if err := provider.Boot(container); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}
