package servlet

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"sync"
)

// ConfigSourceFile defines an instance of a file stream configuration source.
type ConfigSourceFile struct {
	ConfigSourceBase
	path           string
	format         string
	fileSystem     afero.Fs
	decoderFactory *ConfigDecoderFactory
}

// NewConfigSourceFile instantiate a new source that treats a file as
// the origin of the configuration content.
func NewConfigSourceFile(path string, format string, fileSystem afero.Fs, decoderFactory *ConfigDecoderFactory) (*ConfigSourceFile, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("invalid nil 'fileSystem' argument")
	}
	if decoderFactory == nil {
		return nil, fmt.Errorf("invalid nil 'decoderFactory' argument")
	}

	s := &ConfigSourceFile{
		ConfigSourceBase: ConfigSourceBase{
			mutex:   &sync.Mutex{},
			partial: nil,
		},
		path:           path,
		format:         format,
		fileSystem:     fileSystem,
		decoderFactory: decoderFactory,
	}

	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *ConfigSourceFile) load() error {
	file, err := s.fileSystem.OpenFile(s.path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	decoder, err := s.decoderFactory.Create(s.format, file)
	if err != nil {
		_ = file.Close()
		return err
	}
	defer decoder.Close()

	partial, err := decoder.Decode()
	if err != nil {
		return err
	}

	s.mutex.Lock()
	s.partial = partial
	s.mutex.Unlock()

	return nil
}
