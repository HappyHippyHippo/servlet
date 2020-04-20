package config

import (
	"github.com/happyhippyhippo/servlet"
	"github.com/spf13/afero"
)

type provider struct {
	params Parameters
}

// NewProvider will create a new configuration provider instance used to
// register the basic configuration objects in the application container.
func NewProvider(params Parameters) servlet.Provider {
	return &provider{
		params: params,
	}
}

// Register will register the configuration section instances in the
// application container.
func (p provider) Register(container servlet.Container) {
	container.Add(p.params.DecoderFactoryID, func(container servlet.Container) interface{} {
		decoderFactory := NewDecoderFactory()
		decoderFactory.Register(NewYamlDecoderFactoryStrategy())

		return decoderFactory
	})

	container.Add(p.params.SourceFactoryID, func(container servlet.Container) interface{} {
		fileSystem := container.Get(p.params.FileSystemID).(afero.Fs)
		decoderFactory := container.Get(p.params.DecoderFactoryID).(DecoderFactory)

		fileSourceFactoryStrategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)
		obseravbleFileSourceFactoryStrategy, _ := NewObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

		sourceFactory := NewSourceFactory()
		sourceFactory.Register(fileSourceFactoryStrategy)
		sourceFactory.Register(obseravbleFileSourceFactoryStrategy)

		return sourceFactory
	})

	container.Add(p.params.ConfigID, func(container servlet.Container) interface{} {
		config, _ := NewConfig(p.params.ObserveFrequency)
		return config
	})

	container.Add(p.params.LoaderID, func(container servlet.Container) interface{} {
		config := container.Get(p.params.ConfigID).(Config)
		sourceFactory := container.Get(p.params.SourceFactoryID).(SourceFactory)

		loader, _ := NewLoader(config, sourceFactory)
		return loader
	})
}

// Boot will start the configuration config instance by calling the
// configuration loader with the defined provider base entry information.
func (p provider) Boot(container servlet.Container) {
	loader := container.Get(p.params.LoaderID).(Loader)

	loader.Load(p.params.BaseSourceID, p.params.BaseSourcePath, p.params.BaseSourceFormat)
}
