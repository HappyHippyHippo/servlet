package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"io"
	"os"
	"reflect"
	"testing"
)

func Test_NewConfigSourceFile(t *testing.T) {
	t.Run("nil file system adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		decoderFactory := NewConfigDecoderFactory()

		if source, err := NewConfigSourceFile("path", ConfigDecoderFormatYAML, nil, decoderFactory); source != nil {
			defer source.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil decoder factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		if source, err := NewConfigSourceFile("path", ConfigDecoderFormatYAML, fileSystem, nil); source != nil {
			defer source.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'decoderFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when opening the file", func(t *testing.T) {
		path := "path"
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)).Times(1)
		decoderFactory := NewConfigDecoderFactory()

		if source, err := NewConfigSourceFile(path, ConfigDecoderFormatYAML, fileSystem, decoderFactory); source != nil {
			defer source.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when creating the decoder", func(t *testing.T) {
		path := "path"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()

		if source, err := NewConfigSourceFile(path, "invalid_format", fileSystem, decoderFactory); source != nil {
			defer source.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized format type : invalid_format" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when running the decoder", func(t *testing.T) {
		path := "path"
		errorMessage := "error"
		expectedError := fmt.Sprintf("yaml: input error: %s", errorMessage)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			return 0, fmt.Errorf(errorMessage)
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		if source, err := NewConfigSourceFile(path, ConfigDecoderFormatYAML, fileSystem, decoderFactory); source != nil {
			defer source.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("creates the config file source", func(t *testing.T) {
		path := "path"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, "field: value")
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		if source, err := NewConfigSourceFile(path, ConfigDecoderFormatYAML, fileSystem, decoderFactory); source == nil {
			t.Error("didn't returned a valid reference")
		} else {
			defer source.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if source.mutex == nil {
				t.Error("didn't created the access mutex")
			} else if source.path != path {
				t.Error("didn't stored the file path")
			} else if source.format != ConfigDecoderFormatYAML {
				t.Error("didn't stored the file content format")
			} else if source.fileSystem != fileSystem {
				t.Error("didn't stored the file system adapter reference")
			} else if source.decoderFactory != decoderFactory {
				t.Error("didn't stored the decoder factory reference")
			}
		}
	})

	t.Run("store the decoded partial", func(t *testing.T) {
		path := "path"
		field := "field"
		value := "value"
		expected := ConfigPartial{field: value}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		source, _ := NewConfigSourceFile(path, ConfigDecoderFormatYAML, fileSystem, decoderFactory)

		if check := source.partial; !reflect.DeepEqual(check, expected) {
			t.Error("didn't correctly stored the decoded partial")
		}
	})
}
