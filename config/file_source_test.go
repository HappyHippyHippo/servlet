package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewFileSource(t *testing.T) {
	t.Run("should return an error when missing the file system", func(t *testing.T) {
		action := "Creating a new file source without the file system"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "Invalid nil 'fileSystem' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		decoderFactory := NewMockDecoderFactory(ctrl)

		stream, err := NewFileSource(path, format, nil, decoderFactory)

		if stream != nil {
			t.Errorf("%s return an unexpected valid reference to a new file source", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return an error when missing the decoder factory", func(t *testing.T) {
		action := "Creating a new file source without the decoder factory"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "Invalid nil 'decoderFactory' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		stream, err := NewFileSource(path, format, fileSystem, nil)

		if stream != nil {
			t.Errorf("%s return an unexpected valid reference to a new file source", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return the error that may be raised when opening the file", func(t *testing.T) {
		action := "Creating a new file source when erroring opening the file"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expected)).Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)

		stream, err := NewFileSource(path, format, fileSystem, decoderFactory)

		if stream != nil {
			stream.Close()
			t.Errorf("%s return an unexpected valid reference to a new file source", action)
		}

		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return the error that may be raised when creating the decoder", func(t *testing.T) {
		action := "Creating a new file source when erroing when creating the decoder"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(nil, fmt.Errorf(expected)).Times(1)

		stream, err := NewFileSource(path, format, fileSystem, decoderFactory)

		if stream != nil {
			stream.Close()
			t.Errorf("%s return an unexpected valid reference to a new file source", action)
		}

		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return the error that may be raised when running the decoder", func(t *testing.T) {
		action := "Creating a new file source when erroing when running the decoder"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(nil, fmt.Errorf(expected)).Times(1)
		decoder.EXPECT().Close().Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, err := NewFileSource(path, format, fileSystem, decoderFactory)

		if stream != nil {
			stream.Close()
			t.Errorf("%s return an unexpected valid reference to a new file source", action)
		}

		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should correctly create the config file source", func(t *testing.T) {
		action := "Creating a new file source"

		path := "__dummy_path__"
		format := DecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		partial := NewMockPartial(ctrl)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(partial, nil).Times(1)
		decoder.EXPECT().Close().Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, err := NewFileSource(path, format, fileSystem, decoderFactory)

		if stream == nil {
			t.Errorf("%s didn't returned the expected config file source", action)
		}
		defer stream.Close()
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}
	})

	t.Run("should correctly correctly store the decoder partial", func(t *testing.T) {
		action := "Storing the decoder partial"

		path := "__dummy_path__"
		format := DecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		partial := NewMockPartial(ctrl)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(partial, nil).Times(1)
		decoder.EXPECT().Close().Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, _ := NewFileSource(path, format, fileSystem, decoderFactory)

		if check := stream.(*fileSource).partial; check != partial {
			t.Errorf("%s didn't correctly stored the decoder returned partial, stored (%v), expected (%v)", action, check, partial)
		}
	})
}
