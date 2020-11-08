package servlet

import (
	"fmt"
	"io"
	"sort"
)

// LogStreamFile defines a file output log stream.
type LogStreamFile struct {
	LogStreamBase
	writer io.Writer
}

// NewLogStreamFile instantiate a new file stream object that will write logging
// content into a file.
func NewLogStreamFile(writer io.Writer, formatter LogFormatter, channels []string, level LogLevel) (LogStream, error) {
	if formatter == nil {
		return nil, fmt.Errorf("invalid nil 'formatter' argument")
	}
	if writer == nil {
		return nil, fmt.Errorf("invalid nil 'writer' argument")
	}

	s := &LogStreamFile{
		LogStreamBase: LogStreamBase{
			formatter,
			channels,
			level},
		writer: writer}

	sort.Strings(s.channels)

	return s, nil
}

// Close will terminate the stream stored writer instance.
func (s *LogStreamFile) Close() (err error) {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if s.writer != nil {
		switch s.writer.(type) {
		case io.Closer:
			err = s.writer.(io.Closer).Close()
		}
		s.writer = nil
	}
	return err
}

// Signal will process the logging signal request and store the logging request
// into the underlying file if passing the channel and level filtering.
func (s LogStreamFile) Signal(channel string, level LogLevel, message string, context map[string]interface{}) error {
	i := sort.SearchStrings(s.channels, channel)
	if i == len(s.channels) || s.channels[i] != channel {
		return nil
	}
	return s.Broadcast(level, message, context)
}

// Broadcast will process the logging signal request and store the logging
// request into the underlying file if passing the level filtering.
func (s LogStreamFile) Broadcast(level LogLevel, message string, context map[string]interface{}) error {
	if s.level < level {
		return nil
	}

	_, err := fmt.Fprintln(s.writer, s.format(level, message, context))
	return err
}
