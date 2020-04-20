package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewFileSource(t *testing.T) {
	path := "path"
	format := DecoderFormatYAML
	expectedError := "error"

	t.Run("error when missing the file system adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		decoderFactory := NewMockDecoderFactory(ctrl)

		if stream, err := NewFileSource("path", DecoderFormatYAML, nil, decoderFactory); stream != nil {
			defer stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error when missing the decoder factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		if stream, err := NewFileSource("path", DecoderFormatYAML, fileSystem, nil); stream != nil {
			defer stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'decoderFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when opening the file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)).Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)

		if stream, err := NewFileSource(path, DecoderFormatYAML, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when creating the decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(nil, fmt.Errorf(expectedError)).Times(1)

		if stream, err := NewFileSource(path, format, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when running the decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(nil, fmt.Errorf(expectedError)).Times(1)
		decoder.EXPECT().Close().Times(1)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		if stream, err := NewFileSource(path, format, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("creates the config file source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		conf := NewMockPartial(ctrl)
		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(conf, nil).Times(1)
		decoder.EXPECT().Close().Times(1)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		if stream, err := NewFileSource(path, format, fileSystem, decoderFactory); stream == nil {
			t.Errorf("didn't return a valid reference")
		} else {
			defer stream.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})

	t.Run("store the decoded partial", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		conf := NewMockPartial(ctrl)
		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(conf, nil).Times(1)
		decoder.EXPECT().Close().Times(1)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, _ := NewFileSource(path, format, fileSystem, decoderFactory)
		defer stream.Close()

		if check := stream.(*fileSource).partial; check != conf {
			t.Errorf("didn't correctly stored the decoded partial")
		}
	})
}
