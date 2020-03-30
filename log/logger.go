package log

import (
	"fmt"
)

// Logger interface defines the methods of a logger instance.
type Logger interface {
	Close() error
	Signal(channel string, level Level, fields F, message string) error
	Broadcast(level Level, fields F, message string) error
	HasStream(id string) bool
	AddStream(id string, stream Stream) error
	RemoveStream(id string)
	Stream(id string) Stream
}

type logger struct {
	streams map[string]Stream
}

// NewLogger create a new logger instance.
func NewLogger() Logger {
	return &logger{
		streams: map[string]Stream{},
	}
}

// Close will terminate all the logging stream associated to the logger.
func (l *logger) Close() error {
	for id, stream := range l.streams {
		stream.Close()
		delete(l.streams, id)
	}
	return nil
}

// Signal will propagate the channel filtered logging request
// to all stored logging streams.
func (l logger) Signal(channel string, level Level, fields F, message string) error {
	for _, stream := range l.streams {
		if err := stream.Signal(channel, level, fields, message); err != nil {
			return err
		}
	}
	return nil
}

// Broadcast will propagate the logging request to all stored logging streams.
func (l logger) Broadcast(level Level, fields F, message string) error {
	for _, stream := range l.streams {
		if err := stream.Broadcast(level, fields, message); err != nil {
			return err
		}
	}
	return nil
}

// HasStream check if a stream is registed with the requested id.
func (l logger) HasStream(id string) bool {
	_, ok := l.streams[id]
	return ok
}

// AddStream registers a new stream into the logger instance.
func (l *logger) AddStream(id string, stream Stream) error {
	if stream == nil {
		return fmt.Errorf("Invalid nil 'stream' argument")
	}

	if l.HasStream(id) {
		return fmt.Errorf("Duplicate id : %s", id)
	}

	l.streams[id] = stream
	return nil
}

// RemoveStream will remove a registed stream with the requested id
// from the logger.
func (l *logger) RemoveStream(id string) {
	if stream, ok := l.streams[id]; ok {
		stream.Close()
		delete(l.streams, id)
	}
}

// Stream retrieve a stream from the logger that is registed with the
// requested id.
func (l logger) Stream(id string) Stream {
	if stream, ok := l.streams[id]; ok {
		return stream
	}
	return nil
}
