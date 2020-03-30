package config

import (
	"github.com/happyhippyhippo/servlet"
	"github.com/spf13/afero"
)

type provider struct {
	params providerParameters
}

// NewProvider will create a new configuration provider instance.
func NewProvider(parameters ProviderParameters) servlet.Provider {
	if parameters == nil {
		parameters = NewDefaultProviderParameters()
	}

	return &provider{
		params: providerParameters{
			id:               parameters.GetID(),
			fileSystemID:     parameters.GetFileSystemID(),
			loaderID:         parameters.GetLoaderID(),
			sourceFactoryID:  parameters.GetSourceFactoryID(),
			decoderFactoryID: parameters.GetDecoderFactoryID(),
			observeFrequency: parameters.GetObserveFrequency(),
			baseSourceID:     parameters.GetBaseSourceID(),
			baseSourcePath:   parameters.GetBaseSourcePath(),
			baseSourceFormat: parameters.GetBaseSourceFormat(),
		},
	}
}

// Register will register the configuration package instances in the
// application container.
func (p provider) Register(container servlet.Container) {
	container.Add(p.params.decoderFactoryID, func(container servlet.Container) interface{} {
		decoderFactory := NewDecoderFactory()
		decoderFactory.Register(NewYamlDecoderFactoryStrategy())

		return decoderFactory
	})

	container.Add(p.params.sourceFactoryID, func(container servlet.Container) interface{} {
		fileSystem := container.Get(p.params.fileSystemID).(afero.Fs)
		decoderFactory := container.Get(p.params.decoderFactoryID).(DecoderFactory)

		fileSourceFactoryStrategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)
		obseravbleFileSourceFactoryStrategy, _ := NewObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

		sourceFactory := NewSourceFactory()
		sourceFactory.Register(fileSourceFactoryStrategy)
		sourceFactory.Register(obseravbleFileSourceFactoryStrategy)

		return sourceFactory
	})

	container.Add(p.params.id, func(container servlet.Container) interface{} {
		config, _ := NewConfig(p.params.observeFrequency)
		return config
	})

	container.Add(p.params.loaderID, func(container servlet.Container) interface{} {
		config := container.Get(p.params.id).(Config)
		sourceFactory := container.Get(p.params.sourceFactoryID).(SourceFactory)

		loader, _ := NewLoader(config, sourceFactory)
		return loader
	})
}

// Boot will start the configuration package config instance by calling the
// configuration loader with the defined provider base entry information.
func (p provider) Boot(container servlet.Container) {
	loader := container.Get(p.params.loaderID).(Loader)

	loader.Load(p.params.baseSourceID, p.params.baseSourcePath, p.params.baseSourceFormat)
}
