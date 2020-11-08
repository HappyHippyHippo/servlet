package servlet

import (
	"fmt"
	"github.com/spf13/afero"
	"sync"
	"time"
)

// ConfigSourceObservableFile defines an instance of a file stream
// configuration source that will be checked for changes periodically in a
// config defined frequency.
type ConfigSourceObservableFile struct {
	ConfigSourceFile
	timestamp time.Time
}

// NewConfigSourceObservableFile instantiate a new source that treats a file
// as the origin of the configuration content. This file source will be
// periodically checked for changes and loaded if so.
func NewConfigSourceObservableFile(path string, format string, fileSystem afero.Fs, decoderFactory *ConfigDecoderFactory) (*ConfigSourceObservableFile, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("invalid nil 'fileSystem' argument")
	}
	if decoderFactory == nil {
		return nil, fmt.Errorf("invalid nil 'decoderFactory' argument")
	}

	s := &ConfigSourceObservableFile{
		ConfigSourceFile: ConfigSourceFile{
			ConfigSourceBase: ConfigSourceBase{
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
func (s *ConfigSourceObservableFile) Reload() (bool, error) {
	if s == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

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
