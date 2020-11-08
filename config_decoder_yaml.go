package servlet

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

type underlyingConfigDecoderYaml interface {
	Decode(partial interface{}) error
}

// ConfigDecoderYaml defines an instance used to decode s YAML encoded config
// source stream
type ConfigDecoderYaml struct {
	reader  io.Reader
	decoder underlyingConfigDecoderYaml
}

// NewConfigDecoderYaml instantiate a new yaml configuration decoder object
// used to parse a yaml configuration source into a config partial.
func NewConfigDecoderYaml(reader io.Reader) (*ConfigDecoderYaml, error) {
	if reader == nil {
		return nil, fmt.Errorf("invalid nil 'reader' argument")
	}

	return &ConfigDecoderYaml{
		reader:  reader,
		decoder: yaml.NewDecoder(reader),
	}, nil
}

// Close terminate the decoder, closing the associated reader.
func (d *ConfigDecoderYaml) Close() {
	if d == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if d.reader != nil {
		switch d.reader.(type) {
		case io.Closer:
			_ = d.reader.(io.Closer).Close()
		}
		d.reader = nil
	}
}

// Decode parse the associated configuration source reader content
// into a configuration partial.
func (d ConfigDecoderYaml) Decode() (ConfigPartial, error) {
	p := ConfigPartial{}
	if err := d.decoder.Decode(&p); err != nil {
		return nil, err
	}
	return p, nil
}
