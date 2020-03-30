package config

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/spf13/afero"
)

type fileSource struct {
	source
	path           string
	format         string
	fileSystem     afero.Fs
	decoderFactory DecoderFactory
}

// NewFileSource instantiate a new file configuration source.
func NewFileSource(path string, format string, fileSystem afero.Fs, decoderFactory DecoderFactory) (Source, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("Invalid nil 'fileSystem' argument")
	}
	if decoderFactory == nil {
		return nil, fmt.Errorf("Invalid nil 'decoderFactory' argument")
	}

	s := &fileSource{
		source: source{
			mutex:   &sync.RWMutex{},
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

func (s *fileSource) load() error {
	file, err := s.fileSystem.OpenFile(s.path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	decoder, err := s.getDecoder(file)
	if err != nil {
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

func (s fileSource) getDecoder(reader io.Reader) (Decoder, error) {
	return s.decoderFactory.Create(s.format, reader)
}
