package log

import (
	"fmt"

	"github.com/happyhippyhippo/servlet/config"
)

// Loader interface defines the methods of a logging loader instance.
type Loader interface {
	Load(c config.Config) error
}

type loader struct {
	formatterFactory FormatterFactory
	streamFactory    StreamFactory
	logger           Logger
}

// NewLoader create a new logging configuration loader instance.
func NewLoader(formatterFactory FormatterFactory, streamFactory StreamFactory, logger Logger) (Loader, error) {
	if formatterFactory == nil {
		return nil, fmt.Errorf("Invalid nil 'formatterFactory' argument")
	}
	if streamFactory == nil {
		return nil, fmt.Errorf("Invalid nil 'streamFactory' argument")
	}
	if logger == nil {
		return nil, fmt.Errorf("Invalid nil 'logger' argument")
	}

	return &loader{
		formatterFactory: formatterFactory,
		streamFactory:    streamFactory,
		logger:           logger,
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

	entries := c.Get("log.streams").([]interface{})
	for _, entry := range entries {
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
