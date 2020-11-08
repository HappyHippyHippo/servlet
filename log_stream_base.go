package servlet

import (
	"fmt"
	"sort"
)

// LogStreamBase defines the base interaction with a log stream instance.
type LogStreamBase struct {
	formatter LogFormatter
	channels  []string
	level     LogLevel
}

// HasChannel will validate if the stream is listening to a specific
// logging channel.
func (s LogStreamBase) HasChannel(channel string) bool {
	i := sort.SearchStrings(s.channels, channel)
	return i < len(s.channels) && s.channels[i] == channel
}

// ListChannels retrieves the list of channels that the stream is listening.
func (s LogStreamBase) ListChannels() []string {
	return s.channels
}

// AddChannel register a channel to the list of channels that the
// stream is listening.
func (s *LogStreamBase) AddChannel(channel string) {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if !s.HasChannel(channel) {
		s.channels = append(s.channels, channel)
		sort.Strings(s.channels)
	}
}

// RemoveChannel removes a channel from the list of channels that the
// stream is listening.
func (s *LogStreamBase) RemoveChannel(channel string) {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	i := sort.SearchStrings(s.channels, channel)
	if i == len(s.channels) || s.channels[i] != channel {
		return
	}
	s.channels = append(s.channels[:i], s.channels[i+1:]...)
}

// Level retrieves the logging level filter value of the stream.
func (s LogStreamBase) Level() LogLevel {
	return s.level
}

func (s LogStreamBase) format(level LogLevel, message string, context map[string]interface{}) string {
	if s.formatter != nil {
		message = s.formatter.Format(level, message, context)
	}
	return message
}
