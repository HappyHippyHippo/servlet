package log

import (
	"fmt"
	"io"
	"sort"
)

type fileStream struct {
	stream
	writer io.Writer
}

// NewFileStream instantiate a new file stream object that will write logging
// content into a file.
func NewFileStream(writer io.Writer, formatter Formatter, channels []string, level Level) (Stream, error) {
	if formatter == nil {
		return nil, fmt.Errorf("Invalid nil 'formatter' argument")
	}
	if writer == nil {
		return nil, fmt.Errorf("Invalid nil 'writer' argument")
	}

	s := &fileStream{
		stream{
			formatter,
			channels,
			level},
		writer}

	sort.Strings(s.channels)

	return s, nil
}

// Close will terminate the stream stored writer instance.
func (s *fileStream) Close() (err error) {
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
func (s fileStream) Signal(channel string, level Level, message string, fields F) error {
	i := sort.SearchStrings(s.channels, channel)
	if i == len(s.channels) || s.channels[i] != channel {
		return nil
	}
	return s.Broadcast(level, message, fields)
}

// Broadcast will process the logging signal request and store the logging
// request into the underlying file if passing the level filtering.
func (s fileStream) Broadcast(level Level, message string, fields F) error {
	if s.level < level {
		return nil
	}

	_, err := fmt.Fprintln(s.writer, s.format(level, message, fields))
	return err
}
