package log

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewFileStreamFactoryStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	formatterFactory := NewMockFormatterFactory(ctrl)

	t.Run("error when missing file system adapter", func(t *testing.T) {
		if strategy, err := NewFileStreamFactoryStrategy(nil, formatterFactory); strategy != nil {
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error when missing the formatter factory", func(t *testing.T) {
		if strategy, err := NewFileStreamFactoryStrategy(fileSystem, nil); strategy != nil {
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'formatterFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("creates a new file stream factory strategy", func(t *testing.T) {
		if strategy, err := NewFileStreamFactoryStrategy(fileSystem, formatterFactory); strategy == nil {
			t.Errorf("didn't return a valid reference")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_FileStreamFactoryStrategy_Accept(t *testing.T) {
	path := "path"
	format := "format"
	channels := []string{}
	level := DEBUG

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	formatterFactory := NewMockFormatterFactory(ctrl)
	strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

	t.Run("don't accept if less then 4 extra arguments", func(t *testing.T) {
		if strategy.Accept(StreamTypeFile, path, format, channels) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if the first extra argument is not a string", func(t *testing.T) {
		if strategy.Accept(StreamTypeFile, []byte{}, format, channels, level) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if the second extra argument is not a string", func(t *testing.T) {
		if strategy.Accept(StreamTypeFile, path, []byte{}, channels, level) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if the third extra argument is not a list of strings", func(t *testing.T) {
		if strategy.Accept(StreamTypeFile, path, format, []byte{}, level) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if the forth extra argument is not a string", func(t *testing.T) {
		if strategy.Accept(StreamTypeFile, path, format, channels, []byte{}) {
			t.Errorf("returned true")
		}
	})

	t.Run("accept only file type", func(t *testing.T) {
		scenarios := []struct {
			stype    string
			expected bool
		}{
			{ // test file type
				stype:    StreamTypeFile,
				expected: true,
			},
			{ // test non-file format (db)
				stype:    "db",
				expected: false,
			},
		}

		for _, scn := range scenarios {
			if check := strategy.Accept(scn.stype, path, format, channels, level); check != scn.expected {
				t.Errorf("returned (%v) for the type (%s)", check, scn.stype)
			}
		}
	})
}

func Test_FileStreamFactoryStrategy_AcceptConfig(t *testing.T) {
	stype := StreamTypeFile
	path := "path"
	format := FormatterFormatJSON
	channels := []interface{}{"channel.1", "channel.2"}
	level := "debug"

	t.Run("don't accept if type is missing or is not a string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("type").DoAndReturn(func(key string) string {
			panic("invalid convertion")
		}).Times(1)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)
		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if path is missing or is not a string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		gomock.InOrder(
			partial.EXPECT().String("type").Return(stype).Times(1),
			partial.EXPECT().String("path").DoAndReturn(func(key string) string {
				panic("invalid convertion")
			}).Times(1),
		)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)
		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if format is missing or is not a string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		gomock.InOrder(
			partial.EXPECT().String("type").Return(stype).Times(1),
			partial.EXPECT().String("path").Return(path).Times(1),
			partial.EXPECT().String("format").DoAndReturn(func(key string) string {
				panic("invalid convertion")
			}).Times(1),
		)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)
		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if channels is missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		gomock.InOrder(
			partial.EXPECT().String("type").Return(stype).Times(1),
			partial.EXPECT().String("path").Return(path).Times(1),
			partial.EXPECT().String("format").Return(format).Times(1),
		)
		partial.EXPECT().Get("channels").Return(nil).Times(1)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)
		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if channels is not a list of strings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		gomock.InOrder(
			partial.EXPECT().String("type").Return(stype).Times(1),
			partial.EXPECT().String("path").Return(path).Times(1),
			partial.EXPECT().String("format").Return(format).Times(1),
		)
		partial.EXPECT().Get("channels").Return([]int{}).Times(1)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)
		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if level is missing or is not a string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		gomock.InOrder(
			partial.EXPECT().String("type").Return(stype).Times(1),
			partial.EXPECT().String("path").Return(path).Times(1),
			partial.EXPECT().String("format").Return(format).Times(1),
			partial.EXPECT().String("level").DoAndReturn(func(key string) string {
				panic("invalid convertion")
			}).Times(1),
		)
		partial.EXPECT().Get("channels").Return(channels).Times(1)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)
		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if level is unrecognizable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		gomock.InOrder(
			partial.EXPECT().String("type").Return(stype).Times(1),
			partial.EXPECT().String("path").Return(path).Times(1),
			partial.EXPECT().String("format").Return(format).Times(1),
			partial.EXPECT().String("level").Return("invalid").Times(1),
		)
		partial.EXPECT().Get("channels").Return(channels).Times(1)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)
		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("returned true")
		}
	})

	t.Run("accept the config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		gomock.InOrder(
			partial.EXPECT().String("type").Return(stype).Times(1),
			partial.EXPECT().String("path").Return(path).Times(1),
			partial.EXPECT().String("format").Return(format).Times(1),
			partial.EXPECT().String("level").Return(level).Times(1),
		)
		partial.EXPECT().Get("channels").Return(channels).Times(1)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)
		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if !strategy.AcceptConfig(partial) {
			t.Errorf("returned false")
		}
	})
}

func Test_FileStreamFactoryStrategy_Create(t *testing.T) {
	path := "path"
	format := "format"
	channels := []string{}
	level := DEBUG
	expectedError := "dummy_error"

	t.Run("error on opening the file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)).Times(1)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if stream, err := strategy.Create(path, format, channels, level); stream != nil {
			stream.Close()
			t.Errorf("return a valie stream")
		} else if err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on creating the formatter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		formatterFactory := NewMockFormatterFactory(ctrl)
		formatterFactory.EXPECT().Create(format).Return(nil, fmt.Errorf(expectedError)).Times(1)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if stream, err := strategy.Create(path, format, channels, level); stream != nil {
			stream.Close()
			t.Errorf("return a valie stream")
		} else if err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the file stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		formatter := NewMockFormatter(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)
		formatterFactory.EXPECT().Create(format).Return(formatter, nil).Times(1)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if stream, err := strategy.Create(path, format, channels, level); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if stream == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			switch stream.(type) {
			case *fileStream:
			default:
				t.Errorf("didn't return a new file stream")
			}
		}
	})
}

func Test_FileStreamFactoryStrategy_CreateConfig(t *testing.T) {
	path := "path"
	format := "format"
	channels := []interface{}{"channel1"}
	level := "debug"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conf := NewMockPartial(ctrl)
	gomock.InOrder(
		conf.EXPECT().String("path").Return(path).Times(1),
		conf.EXPECT().String("format").Return(format).Times(1),
		conf.EXPECT().Get("channels").Return(channels).Times(1),
		conf.EXPECT().String("level").Return(level).Times(1),
	)

	file := NewMockFile(ctrl)
	fileSystem := NewMockFs(ctrl)
	fileSystem.EXPECT().OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
	formatter := NewMockFormatter(ctrl)
	formatterFactory := NewMockFormatterFactory(ctrl)
	formatterFactory.EXPECT().Create(format).Return(formatter, nil).Times(1)

	strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

	if stream, err := strategy.CreateConfig(conf); err != nil {
		t.Errorf("returned the (%v) error", err)
	} else if stream == nil {
		t.Errorf("didn't returned a valid reference")
	} else {
		switch stream.(type) {
		case *fileStream:
		default:
			t.Errorf("didn't return a new file stream")
		}
	}
}
