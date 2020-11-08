package servlet

import "fmt"

// LogLoader defines the log instantiation and initialization of a new
// log proxy.
type LogLoader struct {
	logger        *Log
	streamFactory *LogStreamFactory
}

// NewLogLoader create a new logging configuration loader instance.
func NewLogLoader(logger *Log, streamFactory *LogStreamFactory) (*LogLoader, error) {
	if logger == nil {
		return nil, fmt.Errorf("invalid nil 'logger' argument")
	}
	if streamFactory == nil {
		return nil, fmt.Errorf("invalid nil 'streamFactory' argument")
	}

	return &LogLoader{
		logger:        logger,
		streamFactory: streamFactory,
	}, nil
}

// Load will parse the configuration and instantiates logging streams
// depending the data on the configuration.
func (l LogLoader) Load(c *Config) (err error) {
	if c == nil {
		return fmt.Errorf("invalid nil 'config' argument")
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	entries := c.Get("log.streams")
	if entries == nil {
		return nil
	}

	for _, entry := range entries.([]interface{}) {
		if err = l.load(entry.(ConfigPartial)); err != nil {
			return err
		}
	}

	return nil
}

func (l LogLoader) load(conf ConfigPartial) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	id := conf.String("id")

	var stream LogStream
	if stream, err = l.streamFactory.CreateConfig(conf); err != nil {
		return err
	}

	if err = l.logger.AddStream(id, stream); err != nil {
		return err
	}

	return nil
}
