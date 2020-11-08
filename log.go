package servlet

import (
	"fmt"
)

// Log defines a logging proxy for all the registered logging streams.
type Log struct {
	streams map[string]LogStream
}

// NewLog create a new logger instance.
func NewLog() *Log {
	return &Log{
		streams: map[string]LogStream{},
	}
}

// Close will terminate all the logging stream associated to the logger.
func (l *Log) Close() error {
	if l == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	for id, stream := range l.streams {
		_ = stream.Close()
		delete(l.streams, id)
	}
	return nil
}

// Signal will propagate the channel filtered logging request
// to all stored logging streams.
func (l Log) Signal(channel string, level LogLevel, message string, context map[string]interface{}) error {
	for _, stream := range l.streams {
		if err := stream.Signal(channel, level, message, context); err != nil {
			return err
		}
	}
	return nil
}

// Broadcast will propagate the logging request to all stored logging streams.
func (l Log) Broadcast(level LogLevel, message string, context map[string]interface{}) error {
	for _, stream := range l.streams {
		if err := stream.Broadcast(level, message, context); err != nil {
			return err
		}
	}
	return nil
}

// HasStream check if a stream is registered with the requested id.
func (l Log) HasStream(id string) bool {
	_, ok := l.streams[id]
	return ok
}

// AddStream registers a new stream into the logger instance.
func (l *Log) AddStream(id string, stream LogStream) error {
	if l == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if stream == nil {
		return fmt.Errorf("invalid nil 'stream' argument")
	}

	if l.HasStream(id) {
		return fmt.Errorf("duplicate id : %s", id)
	}

	l.streams[id] = stream
	return nil
}

// RemoveStream will remove a registered stream with the requested id
// from the logger.
func (l *Log) RemoveStream(id string) {
	if l == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if stream, ok := l.streams[id]; ok {
		_ = stream.Close()
		delete(l.streams, id)
	}
}

// Stream retrieve a stream from the logger that is registered with the
// requested id.
func (l Log) Stream(id string) LogStream {
	if stream, ok := l.streams[id]; ok {
		return stream
	}
	return nil
}
