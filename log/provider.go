package log

import (
	"github.com/happyhippyhippo/servlet"
	"github.com/happyhippyhippo/servlet/config"
	"github.com/spf13/afero"
)

type provider struct {
	params Parameters
}

// NewProvider will create a new logger provider instance.
func NewProvider(params Parameters) servlet.Provider {
	return &provider{
		params: params,
	}
}

// Register will register the logger package instances in the
// application container.
func (p provider) Register(container servlet.Container) {
	container.Add(p.params.FormatterFactoryID, func(container servlet.Container) interface{} {
		formatterFactory := NewFormatterFactory()
		formatterFactory.Register(NewJSONFormatterFactoryStrategy())

		return formatterFactory
	})

	container.Add(p.params.StreamFactoryID, func(container servlet.Container) interface{} {
		fileSystem := container.Get(p.params.FileSystemID).(afero.Fs)
		formatterFactory := container.Get(p.params.FormatterFactoryID).(FormatterFactory)

		fileStreamFactoryStrategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		streamFactory := NewStreamFactory()
		streamFactory.Register(fileStreamFactoryStrategy)

		return streamFactory
	})

	container.Add(p.params.LoggerID, func(container servlet.Container) interface{} {
		return NewLogger()
	})

	container.Add(p.params.LoaderID, func(container servlet.Container) interface{} {
		logger := container.Get(p.params.LoggerID).(Logger)
		streamFactory := container.Get(p.params.StreamFactoryID).(StreamFactory)

		loader, _ := NewLoader(logger, streamFactory)
		return loader
	})
}

// Boot will start the logger package config instance by calling the
// logger loader with the defined provider base entry information.
func (p provider) Boot(container servlet.Container) {
	loader := container.Get(p.params.LoaderID).(Loader)
	config := container.Get(p.params.ConfigID).(config.Config)

	loader.Load(config)
}
