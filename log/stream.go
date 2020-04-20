package log

import "sort"

// F defines the type of data fields  to be passed to the logging process
// that describes a list of extra content to be added to the logging entry.
type F map[string]interface{}

// Stream interface defines the interaction methods with a logging stream.
type Stream interface {
	Close() error
	Signal(channel string, level Level, message string, fields F) error
	Broadcast(level Level, message string, fields F) error
	HasChannel(channel string) bool
	ListChannels() []string
	AddChannel(channel string)
	RemoveChannel(channel string)
	Level() Level
}

type stream struct {
	formatter Formatter
	channels  []string
	level     Level
}

// HasChannel will validate if the stream is listening to a specific
// logging channel.
func (s stream) HasChannel(channel string) bool {
	i := sort.SearchStrings(s.channels, channel)
	return i < len(s.channels) && s.channels[i] == channel
}

// ListChannels retrieves the list of channels that the stream is listening.
func (s stream) ListChannels() []string {
	return s.channels
}

// AddChannel register a channel to the list of channels that the
// stream is listening.
func (s *stream) AddChannel(channel string) {
	if !s.HasChannel(channel) {
		s.channels = append(s.channels, channel)
		sort.Strings(s.channels)
	}
}

// RemoveChannel removes a channel from the list of channels that the
// stream is listening.
func (s *stream) RemoveChannel(channel string) {
	i := sort.SearchStrings(s.channels, channel)
	if i == len(s.channels) || s.channels[i] != channel {
		return
	}
	s.channels = append(s.channels[:i], s.channels[i+1:]...)
}

// Level retrieves the logging level filter value of the stream.
func (s stream) Level() Level {
	return s.level
}

func (s stream) format(level Level, message string, fields F) string {
	if s.formatter != nil {
		message = s.formatter.Format(level, message, fields)
	}
	return message
}
