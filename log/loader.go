package log

import (
	"fmt"

	"github.com/happyhippyhippo/servlet/config"
)

// Loader interface used to define the methods of a logger loader.
// This loader will reads a configuration partial with the 'log.streams'
// path and creates the logging streams described in the entries by the
// usage of the stream factory. After the creation of the streans instances,
// they will be added to a log managing instance.
type Loader interface {
	Load(c config.Config) error
}

type loader struct {
	logger        Logger
	streamFactory StreamFactory
}

// NewLoader create a new logging configuration loader instance.
func NewLoader(logger Logger, streamFactory StreamFactory) (Loader, error) {
	if logger == nil {
		return nil, fmt.Errorf("Invalid nil 'logger' argument")
	}
	if streamFactory == nil {
		return nil, fmt.Errorf("Invalid nil 'streamFactory' argument")
	}

	return &loader{
		logger:        logger,
		streamFactory: streamFactory,
	}, nil
}

// Load will parse the configuration and instantiates logging streams
// depending the data on the configuration.
func (l loader) Load(c config.Config) (err error) {
	if c == nil {
		return fmt.Errorf("Invalid nil 'config' argument")
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error while parsing the list of streams")
		}
	}()

	entries := c.Get("log.streams")
	if entries == nil {
		return nil
	}

	for _, entry := range entries.([]interface{}) {
		if err = l.load(entry.(config.Partial)); err != nil {
			return err
		}
	}

	return nil
}

func (l loader) load(conf config.Partial) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error while parsing the logger stream entry : %v", conf)
		}
	}()

	id := conf.String("id")
	var stream Stream

	if stream, err = l.streamFactory.CreateConfig(conf); err != nil {
		return err
	}

	if err = l.logger.AddStream(id, stream); err != nil {
		return err
	}

	return nil
}
