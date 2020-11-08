package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_NewConfigSourceObservableFile(t *testing.T) {
	t.Run("nil file system adapter", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		if source, err := NewConfigSourceObservableFile(path, format, nil, decoderFactory); source != nil {
			defer source.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil decoder factory", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		if source, err := NewConfigSourceObservableFile(path, format, fileSystem, nil); source != nil {
			defer source.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'decoderFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when retrieving the file info", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		expectedError := "error"

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(nil, fmt.Errorf(expectedError)).Times(1)

		if source, err := NewConfigSourceObservableFile(path, format, fileSystem, decoderFactory); source != nil {
			defer source.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when opening the file", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		expectedError := "error"

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)).Times(1)

		if source, err := NewConfigSourceObservableFile(path, format, fileSystem, decoderFactory); source != nil {
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

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		if source, err := NewConfigSourceObservableFile(path, "invalid_format", fileSystem, decoderFactory); source != nil {
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
		format := ConfigDecoderFormatYAML

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("{"))
			return 1, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		if source, err := NewConfigSourceObservableFile(path, format, fileSystem, decoderFactory); source != nil {
			defer source.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "yaml: line 1: did not find expected node content" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the config observable file source", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		field := "field"
		value := "value"
		expected := ConfigPartial{field: value}

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		if source, err := NewConfigSourceObservableFile(path, format, fileSystem, decoderFactory); source == nil {
			t.Errorf("didn't returned a valid reference")
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
			} else if !reflect.DeepEqual(source.partial, expected) {
				t.Error("didn't loaded the content correctly")
			}
		}
	})

	t.Run("store the decoded partial", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		field := "field"
		value := "value"
		expected := ConfigPartial{field: value}

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		source, _ := NewConfigSourceObservableFile(path, format, fileSystem, decoderFactory)
		defer source.Close()

		if check := source.partial; !reflect.DeepEqual(check, expected) {
			t.Error("didn't correctly stored the decoded partial")
		}
	})
}

func Test_ConfigSourceObservableFile_Reload(t *testing.T) {
	t.Run("nil pointer receiver", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("didn't panic")
			} else {
				switch e := r.(type) {
				case error:
					if e.Error() != "nil pointer receiver" {
						t.Errorf("panic with the (%v) error", e)
					}
				default:
					t.Error("didn't panic with an error")
				}
			}
		}()

		var source *ConfigSourceObservableFile
		_, _ = source.Reload()
	})

	t.Run("error if fail to retrieving the file info", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		field := "field"
		value := "value"
		expectedError := "error"

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		gomock.InOrder(
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
			fileSystem.EXPECT().Stat(path).Return(nil, fmt.Errorf(expectedError)),
		)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		source, _ := NewConfigSourceObservableFile(path, format, fileSystem, decoderFactory)
		defer source.Close()

		if reloaded, err := source.Reload(); reloaded {
			t.Error("flagged that was reloaded")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error if fails to load the file content", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		field := "field"
		value := "value"
		expectedError := "error"

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		gomock.InOrder(
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)),
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 2)),
		)
		fileSystem := NewMockFs(ctrl)
		gomock.InOrder(
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
		)
		gomock.InOrder(
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil),
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)),
		)

		source, _ := NewConfigSourceObservableFile(path, format, fileSystem, decoderFactory)
		defer source.Close()

		if reloaded, err := source.Reload(); reloaded {
			t.Error("flagged that was reloaded")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("prevent reload of a unchanged source", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		field := "field"
		value := "value"

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(2)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(2)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		source, _ := NewConfigSourceObservableFile(path, format, fileSystem, decoderFactory)

		if reloaded, err := source.Reload(); reloaded {
			t.Error("flagged that was reloaded")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("should reload a changed source", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		field := "field"
		value1 := "value1"
		value2 := "value2"
		expected := ConfigPartial{field: value2}

		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file1 := NewMockFile(ctrl)
		file1.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value1))
			return 13, io.EOF
		})
		file1.EXPECT().Close().Times(1)
		file2 := NewMockFile(ctrl)
		file2.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value2))
			return 13, io.EOF
		})
		file2.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		gomock.InOrder(
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)),
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 2)),
		)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(2)
		gomock.InOrder(
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file1, nil),
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file2, nil),
		)

		source, _ := NewConfigSourceObservableFile(path, format, fileSystem, decoderFactory)

		if reloaded, err := source.Reload(); !reloaded {
			t.Error("flagged that was not reloaded")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(expected, source.partial) {
			t.Error("didn't stored the reloaded configuration")
		}
	})
}
