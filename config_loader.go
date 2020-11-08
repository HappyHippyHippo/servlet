package servlet

import "fmt"

// ConfigLoader defines the config instantiation and initialization of a new
// config managing structure.
type ConfigLoader struct {
	config        *Config
	sourceFactory *ConfigSourceFactory
}

// NewConfigLoader instantiate a new configuration loader.
func NewConfigLoader(config *Config, sourceFactory *ConfigSourceFactory) (*ConfigLoader, error) {
	if config == nil {
		return nil, fmt.Errorf("invalid nil 'config' argument")
	}
	if sourceFactory == nil {
		return nil, fmt.Errorf("invalid nil 'sourceFactory' argument")
	}

	return &ConfigLoader{
		config:        config,
		sourceFactory: sourceFactory,
	}, nil
}

// Load loads the configuration from a base config file defined by a
// path and format.
func (l ConfigLoader) Load(id string, path string, format string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error while parsing the list of sources")
		}
	}()

	source, err := l.sourceFactory.Create(ConfigSourceTypeFile, path, format)
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
		if err = l.loadSource(conf.(ConfigPartial)); err != nil {
			return err
		}
	}

	return nil
}

func (l ConfigLoader) loadSource(conf ConfigPartial) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	id := conf.String("id")
	priority := conf.Int("priority")

	var source ConfigSource
	if source, err = l.sourceFactory.CreateConfig(conf); err != nil {
		return err
	}

	if err = l.config.AddSource(id, priority, source); err != nil {
		return err
	}

	return nil
}
