package servlet

import (
	"fmt"
	"github.com/spf13/afero"
)

// ConfigProvider defines the default configuration provider to be used on
// the application initialization to register the configuration services.
type ConfigProvider struct {
	params *ConfigProviderParams
}

// NewConfigProvider will create a new configuration provider instance used to
// register the basic configuration objects in the application container.
func NewConfigProvider(params *ConfigProviderParams) *ConfigProvider {
	if params == nil {
		params = NewConfigProviderParams()
	}

	return &ConfigProvider{
		params: params,
	}
}

// Register will register the configuration section instances in the
// application container.
func (p ConfigProvider) Register(container *AppContainer) error {
	if container == nil {
		return fmt.Errorf("invalid nil 'container' argument")
	}

	_ = container.Add(p.params.DecoderFactoryStrategyYamlID, func(container *AppContainer) (interface{}, error) {
		return NewConfigDecoderFactoryStrategyYaml(), nil
	})

	_ = container.Add(p.params.DecoderFactoryID, func(container *AppContainer) (interface{}, error) {
		return NewConfigDecoderFactory(), nil
	})

	_ = container.Add(p.params.SourceFactoryStrategyFileID, func(container *AppContainer) (strategy interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()

		fileSystem, err := container.Get(p.params.FileSystemID)
		if err != nil {
			return nil, err
		}

		decoderFactory, err := container.Get(p.params.DecoderFactoryID)
		if err != nil {
			return nil, err
		}

		return NewConfigSourceFactoryStrategyFile(fileSystem.(afero.Fs), decoderFactory.(*ConfigDecoderFactory))
	})

	_ = container.Add(p.params.SourceFactoryStrategyObservableFileID, func(container *AppContainer) (strategy interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()

		fileSystem, err := container.Get(p.params.FileSystemID)
		if err != nil {
			return nil, err
		}

		decoderFactory, err := container.Get(p.params.DecoderFactoryID)
		if err != nil {
			return nil, err
		}

		return NewConfigSourceFactoryStrategyObservableFile(fileSystem.(afero.Fs), decoderFactory.(*ConfigDecoderFactory))
	})

	_ = container.Add(p.params.SourceFactoryStrategyEnvironmentID, func(container *AppContainer) (interface{}, error) {
		return NewConfigSourceFactoryStrategyEnvironment()
	})

	_ = container.Add(p.params.SourceFactoryID, func(container *AppContainer) (obj interface{}, err error) {
		return NewConfigSourceFactory(), nil
	})

	_ = container.Add(p.params.ConfigID, func(container *AppContainer) (interface{}, error) {
		return NewConfig(p.params.ObserveFrequency)
	})

	_ = container.Add(p.params.LoaderID, func(container *AppContainer) (obj interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()

		config, err := container.Get(p.params.ConfigID)
		if err != nil {
			return nil, err
		}

		sourceFactory, err := container.Get(p.params.SourceFactoryID)
		if err != nil {
			return nil, err
		}

		return NewConfigLoader(config.(*Config), sourceFactory.(*ConfigSourceFactory))
	})

	return nil
}

// Boot will start the configuration config instance by calling the
// configuration loader with the defined provider base entry information.
func (p ConfigProvider) Boot(container *AppContainer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	{
		factory, err := container.Get(p.params.DecoderFactoryID)
		if err != nil {
			return err
		}

		{
			strategy, err := container.Get(p.params.DecoderFactoryStrategyYamlID)
			if err != nil {
				return err
			}

			_ = factory.(*ConfigDecoderFactory).Register(strategy.(ConfigDecoderFactoryStrategy))
		}
	}

	{
		factory, err := container.Get(p.params.SourceFactoryID)
		if err != nil {
			return err
		}

		{
			strategy, err := container.Get(p.params.SourceFactoryStrategyFileID)
			if err != nil {
				return err
			}

			_ = factory.(*ConfigSourceFactory).Register(strategy.(ConfigSourceFactoryStrategy))
		}

		{
			strategy, err := container.Get(p.params.SourceFactoryStrategyObservableFileID)
			if err != nil {
				return err
			}

			_ = factory.(*ConfigSourceFactory).Register(strategy.(ConfigSourceFactoryStrategy))
		}

		{
			strategy, err := container.Get(p.params.SourceFactoryStrategyEnvironmentID)
			if err != nil {
				return err
			}

			_ = factory.(*ConfigSourceFactory).Register(strategy.(ConfigSourceFactoryStrategy))
		}
	}

	if p.params.EntrySourceActive {
		loader, err := container.Get(p.params.LoaderID)
		if err != nil {
			return err
		}

		return loader.(*ConfigLoader).Load(p.params.EntrySourceID, p.params.EntrySourcePath, p.params.EntrySourceFormat)
	}

	return nil
}
