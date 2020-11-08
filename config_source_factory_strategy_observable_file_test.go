package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_NewConfigSourceFactoryStrategyObservableFile(t *testing.T) {
	t.Run("nil file system adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		decoderFactory := NewConfigDecoderFactory()

		if strategy, err := NewConfigSourceFactoryStrategyObservableFile(nil, decoderFactory); strategy != nil {
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

		if strategy, err := NewConfigSourceFactoryStrategyObservableFile(fileSystem, nil); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'decoderFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new file source factory strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()

		if strategy, err := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if strategy == nil {
			t.Error("didn't returned a valid reference")
		} else if strategy.fileSystem != fileSystem {
			t.Error("didn't stored the file system adapter reference")
		} else if strategy.decoderFactory != decoderFactory {
			t.Error("didn't stored the decoder factory reference")
		}
	})
}

func Test_ConfigSourceFactoryStrategyObservableFile_Accept(t *testing.T) {
	t.Run("don't accept if at least 2 extra arguments are passed", func(t *testing.T) {
		sourceType := ConfigSourceTypeObservableFile
		path := "path"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		if strategy.Accept(sourceType, path) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if the path is not a string", func(t *testing.T) {
		sourceType := ConfigSourceTypeObservableFile
		format := ConfigDecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		if strategy.Accept(sourceType, 1, format) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if the format is not a string", func(t *testing.T) {
		sourceType := ConfigSourceTypeObservableFile
		path := "path"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		if strategy.Accept(sourceType, path, 1) {
			t.Error("returned true")
		}
	})

	t.Run("accept only file type", func(t *testing.T) {
		scenarios := []struct {
			sourceType string
			expected   bool
		}{
			{ // test file type
				sourceType: ConfigSourceTypeObservableFile,
				expected:   true,
			},
			{ // test non-file type (file)
				sourceType: ConfigSourceTypeFile,
				expected:   false,
			},
		}

		for _, scn := range scenarios {
			path := "path"
			format := ConfigDecoderFormatYAML

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fileSystem := NewMockFs(ctrl)
			decoderFactory := NewConfigDecoderFactory()
			strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

			if check := strategy.Accept(scn.sourceType, path, format); check != scn.expected {
				t.Errorf("for the type (%s), returned (%v)", scn.sourceType, check)
			}
		}
	})
}

func Test_ConfigSourceFactoryStrategyObservableFile_AcceptConfig(t *testing.T) {
	t.Run("don't accept if type is missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		partial := ConfigPartial{}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if type is not a string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		partial := ConfigPartial{"type": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if path is missing", func(t *testing.T) {
		sourceType := ConfigSourceTypeObservableFile

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		partial := ConfigPartial{"type": sourceType}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if path is not a string", func(t *testing.T) {
		sourceType := ConfigSourceTypeObservableFile

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		partial := ConfigPartial{"type": sourceType, "path": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if format is missing", func(t *testing.T) {
		sourceType := ConfigSourceTypeObservableFile
		path := "path"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		partial := ConfigPartial{"type": sourceType, "path": path}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if format is not a string", func(t *testing.T) {
		sourceType := ConfigSourceTypeObservableFile
		path := "path"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		partial := ConfigPartial{"type": sourceType, "path": path, "format": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if invalid type", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		partial := ConfigPartial{"type": ConfigSourceTypeFile, "path": path, "format": format}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("accept config", func(t *testing.T) {
		sourceType := ConfigSourceTypeObservableFile
		path := "path"
		format := ConfigDecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		partial := ConfigPartial{"type": sourceType, "path": path, "format": format}
		if !strategy.AcceptConfig(partial) {
			t.Error("returned false")
		}
	})
}

func Test_ConfigSourceFactoryStrategyObservableFile_Create(t *testing.T) {
	t.Run("non-string path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		if source, err := strategy.Create(123, "format"); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("non-string format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyFile(fileSystem, decoderFactory)

		if source, err := strategy.Create("path", 123); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the file source", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
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
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		if source, err := strategy.Create(path, format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if source == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch s := source.(type) {
			case *ConfigSourceObservableFile:
				if !reflect.DeepEqual(s.partial, expected) {
					t.Error("didn't loaded the content correctly")
				}
			default:
				t.Error("didn't returned a new file source")
			}
		}
	})
}

func Test_ConfigSourceFactoryStrategyObservableFile_CreateConfig(t *testing.T) {
	t.Run("non-string path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		conf := ConfigPartial{"path": 123, "format": "format"}
		if source, err := strategy.CreateConfig(conf); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("non-string format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigSourceFactoryStrategyFile(fileSystem, decoderFactory)

		conf := ConfigPartial{"path": "path", "format": 123}
		if source, err := strategy.CreateConfig(conf); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the file source", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
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
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())
		strategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)

		conf := ConfigPartial{"path": path, "format": format}

		if source, err := strategy.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if source == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch s := source.(type) {
			case *ConfigSourceObservableFile:
				if !reflect.DeepEqual(s.partial, expected) {
					t.Error("didn't loaded the content correctly")
				}
			default:
				t.Error("didn't returned a new file source")
			}
		}
	})
}
