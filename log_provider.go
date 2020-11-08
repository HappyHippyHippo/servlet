package servlet

import (
	"fmt"
	"github.com/spf13/afero"
)

// LogProvider defines the default logging provider to be used on
// the application initialization to register the logging services.
type LogProvider struct {
	params *LogProviderParams
}

// NewLogProvider will create a new logger provider instance.
func NewLogProvider(params *LogProviderParams) *LogProvider {
	if params == nil {
		params = NewLogProviderParams()
	}

	return &LogProvider{
		params: params,
	}
}

// Register will register the logger package instances in the
// application container.
func (p LogProvider) Register(container *AppContainer) error {
	if container == nil {
		return fmt.Errorf("invalid nil 'container' argument")
	}

	_ = container.Add(p.params.FormatterFactoryStrategyJSONID, func(container *AppContainer) (interface{}, error) {
		return NewLogFormatterFactoryStrategyJSON(), nil
	})

	_ = container.Add(p.params.FormatterFactoryID, func(container *AppContainer) (interface{}, error) {
		return NewLogFormatterFactory(), nil
	})

	_ = container.Add(p.params.StreamFactoryStrategyFileID, func(container *AppContainer) (strategy interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()

		fileSystem, err := container.Get(p.params.FileSystemID)
		if err != nil {
			return nil, err
		}

		formatterFactory, err := container.Get(p.params.FormatterFactoryID)
		if err != nil {
			return nil, err
		}

		return NewLogStreamFactoryStrategyFile(fileSystem.(afero.Fs), formatterFactory.(*LogFormatterFactory))
	})

	_ = container.Add(p.params.StreamFactoryID, func(container *AppContainer) (obj interface{}, err error) {
		return NewLogStreamFactory(), nil
	})

	_ = container.Add(p.params.LoggerID, func(container *AppContainer) (interface{}, error) {
		return NewLog(), nil
	})

	_ = container.Add(p.params.LoaderID, func(container *AppContainer) (obj interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()

		logger, err := container.Get(p.params.LoggerID)
		if err != nil {
			return nil, err
		}

		streamFactory, err := container.Get(p.params.StreamFactoryID)
		if err != nil {
			return nil, err
		}

		return NewLogLoader(logger.(*Log), streamFactory.(*LogStreamFactory))
	})

	return nil
}

// Boot will start the logger package config instance by calling the
// logger loader with the defined provider base entry information.
func (p LogProvider) Boot(container *AppContainer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	{
		factory, err := container.Get(p.params.FormatterFactoryID)
		if err != nil {
			return err
		}

		{
			strategy, err := container.Get(p.params.FormatterFactoryStrategyJSONID)
			if err != nil {
				return err
			}

			_ = factory.(*LogFormatterFactory).Register(strategy.(LogFormatterFactoryStrategy))
		}
	}

	{
		factory, err := container.Get(p.params.StreamFactoryID)
		if err != nil {
			return err
		}

		{
			strategy, err := container.Get(p.params.StreamFactoryStrategyFileID)
			if err != nil {
				return err
			}

			_ = factory.(*LogStreamFactory).Register(strategy.(LogStreamFactoryStrategy))
		}
	}

	loader, err := container.Get(p.params.LoaderID)
	if err != nil {
		return err
	}

	config, err := container.Get(p.params.ConfigID)
	if err != nil {
		return err
	}

	return loader.(*LogLoader).Load(config.(*Config))
}
