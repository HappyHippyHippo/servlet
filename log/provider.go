package log

import (
	"github.com/happyhippyhippo/servlet"
	"github.com/happyhippyhippo/servlet/config"
	"github.com/spf13/afero"
)

type provider struct {
	params providerParameters
}

// NewProvider will create a new logger provider instance.
func NewProvider(parameters ProviderParameters) servlet.Provider {
	if parameters == nil {
		parameters = NewDefaultProviderParameters()
	}

	return &provider{
		params: providerParameters{
			id:                 parameters.GetID(),
			fileSystemID:       parameters.GetFileSystemID(),
			configID:           parameters.GetConfigID(),
			formatterFactoryID: parameters.GetFormatterFactoryID(),
			streamFactoryID:    parameters.GetStreamFactoryID(),
			loaderID:           parameters.GetLoaderID(),
		},
	}
}

// Register will register the logger package instances in the
// application container.
func (p provider) Register(container servlet.Container) {
	container.Add(p.params.formatterFactoryID, func(container servlet.Container) interface{} {
		formatterFactory := NewFormatterFactory()
		formatterFactory.Register(NewJSONFormatterFactoryStrategy())

		return formatterFactory
	})

	container.Add(p.params.streamFactoryID, func(container servlet.Container) interface{} {
		fileSystem := container.Get(p.params.fileSystemID).(afero.Fs)
		formatterFactory := container.Get(p.params.formatterFactoryID).(FormatterFactory)

		fileStreamFactoryStrategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		streamFactory := NewStreamFactory()
		streamFactory.Register(fileStreamFactoryStrategy)

		return streamFactory
	})

	container.Add(p.params.id, func(container servlet.Container) interface{} {
		return NewLogger()
	})

	container.Add(p.params.loaderID, func(container servlet.Container) interface{} {
		formatterFactory := container.Get(p.params.formatterFactoryID).(FormatterFactory)
		streamFactory := container.Get(p.params.streamFactoryID).(StreamFactory)
		logger := container.Get(p.params.id).(Logger)

		loader, _ := NewLoader(formatterFactory, streamFactory, logger)
		return loader
	})
}

// Boot will start the logger package config instance by calling the
// logger loader with the defined provider base entry information.
func (p provider) Boot(container servlet.Container) {
	loader := container.Get(p.params.loaderID).(Loader)
	config := container.Get(p.params.configID).(config.Config)

	loader.Load(config)
}
