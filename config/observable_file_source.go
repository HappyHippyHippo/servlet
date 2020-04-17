package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/afero"
)

type observableFileSource struct {
	fileSource
	timestamp time.Time
}

// NewObservableFileSource instantiate a new source that treats a file as
// the origin of the configuration content. This file source will be periodicaly
// checked for changes and loaded if so.
func NewObservableFileSource(path string, format string, fileSystem afero.Fs, decoderFactory DecoderFactory) (ObservableSource, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("Invalid nil 'fileSystem' argument")
	}
	if decoderFactory == nil {
		return nil, fmt.Errorf("Invalid nil 'decoderFactory' argument")
	}

	s := &observableFileSource{
		fileSource: fileSource{
			source: source{
				mutex:   &sync.RWMutex{},
				partial: nil,
			},
			path:           path,
			format:         format,
			fileSystem:     fileSystem,
			decoderFactory: decoderFactory,
		},
		timestamp: time.Unix(0, 0),
	}

	if _, err := s.Reload(); err != nil {
		return nil, err
	}

	return s, nil
}

// Reload will check if the source has been updated, and, if so, reload the
// source configuration partial content.
func (s *observableFileSource) Reload() (bool, error) {
	fileInfo, err := s.fileSystem.Stat(s.path)
	if err != nil {
		return false, err
	}

	info := fileInfo.ModTime()
	if s.timestamp.Equal(time.Unix(0, 0)) || s.timestamp.Before(info) {
		if err := s.load(); err != nil {
			return false, err
		}

		s.mutex.Lock()
		s.timestamp = info
		s.mutex.Unlock()

		return true, nil
	}

	return false, nil
}
