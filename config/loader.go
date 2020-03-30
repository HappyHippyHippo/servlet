package config

import "fmt"

// Loader interface used to define the methods of a configuration loader.
type Loader interface {
	Load(id string, path string, format string) error
}

type loader struct {
	config        Config
	sourceFactory SourceFactory
}

// NewLoader instantiate a new configuration loader.
func NewLoader(config Config, sourceFactory SourceFactory) (Loader, error) {
	if config == nil {
		return nil, fmt.Errorf("Invalid nil 'config' argument")
	}
	if sourceFactory == nil {
		return nil, fmt.Errorf("Invalid nil 'sourceFactory' argument")
	}

	return &loader{
		config:        config,
		sourceFactory: sourceFactory,
	}, nil
}

// Load loads the configuration from a base config file.
func (l loader) Load(id string, path string, format string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error while parsing the list of sources")
		}
	}()

	source, err := l.sourceFactory.Create(SourceTypeFile, path, format)
	if err != nil {
		return err
	}
	if err = l.config.AddSource(id, 0, source); err != nil {
		return err
	}

	entries := l.config.Get("config.sources")
	if entries == nil {
		return nil
	}

	for _, conf := range entries.([]interface{}) {
		if err = l.loadSource(conf.(Partial)); err != nil {
			return err
		}
	}

	return nil
}

func (l loader) loadSource(conf Partial) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error while parsing the config entry : %v", conf)
		}
	}()

	id := conf.String("id")
	priority := conf.Int("priority")
	var source Source

	if source, err = l.sourceFactory.CreateConfig(conf); err != nil {
		return err
	}

	if err = l.config.AddSource(id, priority, source); err != nil {
		return err
	}

	return nil
}
