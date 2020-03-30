package log

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewFileStreamFactoryStrategy(t *testing.T) {
	t.Run("should return nil when missing file system adapter", func(t *testing.T) {
		action := "Creating a file stream factory strategy without a file system adapter reference"

		expected := "Invalid nil 'fileSystem' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, err := NewFileStreamFactoryStrategy(nil, formatterFactory)

		if strategy != nil {
			t.Errorf("%s returned a valid file stream factory strategy, expected nil", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return nil when missing the formatter factory", func(t *testing.T) {
		action := "Creating a file stream factory strategy without a formatter factory reference"

		expected := "Invalid nil 'formatterFactory' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		strategy, err := NewFileStreamFactoryStrategy(fileSystem, nil)

		if strategy != nil {
			t.Errorf("%s returned a valid file stream factory strategy, expected nil", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("creates a new file stream factory strategy", func(t *testing.T) {
		action := "Creating a file stream factory strategy"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, err := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy == nil {
			t.Errorf("%s didn't return a valid reference to a new file stream factory strategy", action)
		}
		if err != nil {
			t.Errorf("%s return a unexpected error : %v", action, err)
		}
	})
}

func Test_FileStreamFactoryStrategy_Accept(t *testing.T) {
	t.Run("should accept only file type", func(t *testing.T) {
		action := "Checking the accepting type"

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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fileSystem := NewMockFs(ctrl)
			formatterFactory := NewMockFormatterFactory(ctrl)

			strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

			path := "__path__"
			format := "__format__"
			channels := []string{}
			level := DEBUG

			if check := strategy.Accept(scn.stype, path, format, channels, level); check != scn.expected {
				t.Errorf("%s didn't returned the expected (%v) for the type (%s), returned (%v)", action, scn.expected, scn.stype, check)
			}
		}
	})

	t.Run("should not accept if less then 4 extra arguments are passed", func(t *testing.T) {
		action := "Checking the acceptance with less then 4 extra arguments"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		path := "__path__"
		format := "__format__"
		channels := []string{}

		if strategy.Accept(StreamTypeFile, path, format, channels) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if the first extra argument is not a string", func(t *testing.T) {
		action := "Checking the acceptance when the first extra argument not a string"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		path := []byte{}
		format := "__format__"
		channels := []string{}
		level := DEBUG

		if strategy.Accept(StreamTypeFile, path, format, channels, level) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if the second extra argument is not a string", func(t *testing.T) {
		action := "Checking the acceptance when the second extra argument not a string"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		path := "__path__"
		format := 1
		channels := []string{}
		level := DEBUG

		if strategy.Accept(StreamTypeFile, path, format, channels, level) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if the third extra argument is not a array of strings", func(t *testing.T) {
		action := "Checking the acceptance when the third extra argument not a array of strings"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		path := "__path__"
		format := "__format__"
		channels := []int{}
		level := DEBUG

		if strategy.Accept(StreamTypeFile, path, format, channels, level) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if the forth extra argument is not a logging level", func(t *testing.T) {
		action := "Checking the acceptance when the forth extra argument not a logging level"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		path := "__path__"
		format := "__format__"
		channels := []string{}
		level := "debug"

		if strategy.Accept(StreamTypeFile, path, format, channels, level) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})
}

func Test_FileStreamFactoryStrategy_AcceptConfig(t *testing.T) {
	t.Run("should not accept if there is not a type config entry", func(t *testing.T) {
		action := "Checking the acceptance of a config without a type field"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)
		conf.EXPECT().String("type").DoAndReturn(func(key string) string {
			panic("invalid convertion")
		})

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(conf) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if there is not a path config entry", func(t *testing.T) {
		action := "Checking the acceptance of a config without a path field"

		stype := StreamTypeFile

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)
		gomock.InOrder(
			conf.EXPECT().String("type").Return(stype).Times(1),
			conf.EXPECT().String("path").DoAndReturn(func(key string) string {
				panic("invalid convertion")
			}),
		)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(conf) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if there is not a format config entry", func(t *testing.T) {
		action := "Checking the acceptance of a config without a format field"

		stype := StreamTypeFile
		path := "__path__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)
		gomock.InOrder(
			conf.EXPECT().String("type").Return(stype).Times(1),
			conf.EXPECT().String("path").Return(path).Times(1),
			conf.EXPECT().String("format").DoAndReturn(func(key string) string {
				panic("invalid convertion")
			}),
		)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(conf) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if there is not a channels config entry", func(t *testing.T) {
		action := "Checking the acceptance of a config without a channels field"

		stype := StreamTypeFile
		path := "__path__"
		format := "__format__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)
		gomock.InOrder(
			conf.EXPECT().String("type").Return(stype).Times(1),
			conf.EXPECT().String("path").Return(path).Times(1),
			conf.EXPECT().String("format").Return(format).Times(1),
			conf.EXPECT().Get("channels").DoAndReturn(func(key string) string {
				panic("invalid convertion")
			}),
		)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(conf) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if the channels config entry is not a list of strings", func(t *testing.T) {
		action := "Checking the acceptance of a config with a channels field not being a list of strings"

		stype := StreamTypeFile
		path := "__path__"
		format := "__format__"
		channels := []int{1, 2}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)
		gomock.InOrder(
			conf.EXPECT().String("type").Return(stype).Times(1),
			conf.EXPECT().String("path").Return(path).Times(1),
			conf.EXPECT().String("format").Return(format).Times(1),
			conf.EXPECT().Get("channels").Return(channels),
		)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(conf) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if there is not a level config entry", func(t *testing.T) {
		action := "Checking the acceptance of a config without a level field"

		stype := StreamTypeFile
		path := "__path__"
		format := "__format__"
		channels := []interface{}{"channel1"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)
		gomock.InOrder(
			conf.EXPECT().String("type").Return(stype).Times(1),
			conf.EXPECT().String("path").Return(path).Times(1),
			conf.EXPECT().String("format").Return(format).Times(1),
			conf.EXPECT().Get("channels").Return(channels).Times(1),
			conf.EXPECT().String("level").DoAndReturn(func(key string) string {
				panic("invalid convertion")
			}),
		)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(conf) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if the level config entry is not recognized", func(t *testing.T) {
		action := "Checking the acceptance of a config unrecognized level field"

		stype := StreamTypeFile
		path := "__path__"
		format := "__format__"
		channels := []interface{}{"channel1"}
		level := "__dummy_level__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)
		gomock.InOrder(
			conf.EXPECT().String("type").Return(stype).Times(1),
			conf.EXPECT().String("path").Return(path).Times(1),
			conf.EXPECT().String("format").Return(format).Times(1),
			conf.EXPECT().Get("channels").Return(channels).Times(1),
			conf.EXPECT().String("level").Return(level).Times(1),
		)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if strategy.AcceptConfig(conf) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should accept a valid config entry", func(t *testing.T) {
		action := "Checking the acceptance of a config"

		stype := StreamTypeFile
		path := "__path__"
		format := "__format__"
		channels := []interface{}{"channel1"}
		level := "debug"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)
		gomock.InOrder(
			conf.EXPECT().String("type").Return(stype).Times(1),
			conf.EXPECT().String("path").Return(path).Times(1),
			conf.EXPECT().String("format").Return(format).Times(1),
			conf.EXPECT().Get("channels").Return(channels).Times(1),
			conf.EXPECT().String("level").Return(level).Times(1),
		)

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if !strategy.AcceptConfig(conf) {
			t.Errorf("%s didn't returned the expected true value", action)
		}
	})
}

func Test_FileStreamFactoryStrategy_Create(t *testing.T) {
	t.Run("should return the error that may occure when opening the file", func(t *testing.T) {
		action := "Creating a new file stream when erroring while opening the file"

		path := "__path__"
		format := "__format__"
		channels := []string{}
		level := DEBUG
		expectedError := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)).Times(1)
		formatterFactory := NewMockFormatterFactory(ctrl)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		stream, err := strategy.Create(path, format, channels, level)
		if stream != nil {
			t.Errorf("%s return a unexpected stream", action)
		}
		if err == nil {
			t.Errorf("%s didn't returned the expected error", action)
		} else {
			if check := err.Error(); check != expectedError {
				t.Errorf("%s returned the error (%s), when expected (%s)", action, check, expectedError)
			}
		}
	})

	t.Run("should return the error that may occure when creating the formatter", func(t *testing.T) {
		action := "Creating a new file stream when erroring while creating the formatter"

		path := "__path__"
		format := "__format__"
		channels := []string{}
		level := DEBUG
		expectedError := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		formatterFactory := NewMockFormatterFactory(ctrl)
		formatterFactory.EXPECT().Create(format).Return(nil, fmt.Errorf(expectedError)).Times(1)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		stream, err := strategy.Create(path, format, channels, level)
		if stream != nil {
			t.Errorf("%s return a unexpected stream", action)
		}
		if err == nil {
			t.Errorf("%s didn't returned the expected error", action)
		} else {
			if check := err.Error(); check != expectedError {
				t.Errorf("%s returned the error (%s), when expected (%s)", action, check, expectedError)
			}
		}
	})

	t.Run("should create the requested file stream", func(t *testing.T) {
		action := "Creating a new file stream"

		path := "__path__"
		format := "__format__"
		channels := []string{}
		level := DEBUG

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		formatter := NewMockFormatter(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		formatterFactory := NewMockFormatterFactory(ctrl)
		formatterFactory.EXPECT().Create(format).Return(formatter, nil).Times(1)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		stream, err := strategy.Create(path, format, channels, level)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if stream == nil {
			t.Errorf("%s didn't returned a valid file stream reference", action)
		} else {
			switch stream.(type) {
			case *fileStream:
			default:
				t.Errorf("%s didn't return a valid reference to a new file stream reference", action)
			}
		}
	})
}

func Test_FileStreamFactoryStrategy_CreateConfig(t *testing.T) {
	t.Run("should create the requested file stream", func(t *testing.T) {
		action := "Creating a new file stream"

		path := "__path__"
		format := "__format__"
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
		formatter := NewMockFormatter(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		formatterFactory := NewMockFormatterFactory(ctrl)
		formatterFactory.EXPECT().Create(format).Return(formatter, nil).Times(1)

		strategy, _ := NewFileStreamFactoryStrategy(fileSystem, formatterFactory)

		stream, err := strategy.CreateConfig(conf)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if stream == nil {
			t.Errorf("%s didn't returned a valid file stream reference", action)
		} else {
			switch stream.(type) {
			case *fileStream:
			default:
				t.Errorf("%s didn't return a valid reference to a new file stream reference", action)
			}
		}
	})
}
