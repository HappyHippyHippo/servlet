package config

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

type underlyingYamlDecoder interface {
	Decode(partial interface{}) error
}

type yamlDecoder struct {
	reader  io.Reader
	decoder underlyingYamlDecoder
}

// NewYamlDecoder instantiate a new yaml configuration decoder object used to
// parse a yaml configuration source into a config partial.
func NewYamlDecoder(reader io.Reader) (Decoder, error) {
	if reader == nil {
		return nil, fmt.Errorf("Invalid nil 'reader' argument")
	}

	return &yamlDecoder{
		reader:  reader,
		decoder: yaml.NewDecoder(reader),
	}, nil
}

// Close terminate the decoder, closing the associated reader.
func (d *yamlDecoder) Close() (err error) {
	if d.reader != nil {
		switch d.reader.(type) {
		case io.Closer:
			err = d.reader.(io.Closer).Close()
		}
		d.reader = nil
	}
	return err
}

// Decode parse the associated configuration source reader content
// into a configuration partial.
func (d yamlDecoder) Decode() (Partial, error) {
	p := partial{}
	if err := d.decoder.Decode(&p); err != nil {
		return nil, err
	}
	return p, nil
}
