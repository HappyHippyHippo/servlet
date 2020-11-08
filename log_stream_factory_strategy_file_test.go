package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"os"
	"strings"
	"testing"
)

func Test_NewLogStreamFactoryStrategyFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	formatterFactory := NewLogFormatterFactory()

	t.Run("nil file system adapter", func(t *testing.T) {
		if strategy, err := NewLogStreamFactoryStrategyFile(nil, formatterFactory); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil formatter factory", func(t *testing.T) {
		if strategy, err := NewLogStreamFactoryStrategyFile(fileSystem, nil); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'formatterFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new file stream factory strategy", func(t *testing.T) {
		if strategy, err := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory); strategy == nil {
			t.Errorf("didn't returned a valid reference")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_LogStreamFactoryStrategyFile_Accept(t *testing.T) {
	path := "path"
	format := "format"
	var channels []string
	level := DEBUG

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	formatterFactory := NewLogFormatterFactory()
	strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

	t.Run("don't accept if less then 4 extra arguments", func(t *testing.T) {
		if strategy.Accept(LogStreamTypeFile, path, format, channels) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if the first extra argument is not a string", func(t *testing.T) {
		if strategy.Accept(LogStreamTypeFile, []byte{}, format, channels, level) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if the second extra argument is not a string", func(t *testing.T) {
		if strategy.Accept(LogStreamTypeFile, path, []byte{}, channels, level) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if the third extra argument is not a list of strings", func(t *testing.T) {
		if strategy.Accept(LogStreamTypeFile, path, format, []byte{}, level) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if the forth extra argument is not a string", func(t *testing.T) {
		if strategy.Accept(LogStreamTypeFile, path, format, channels, []byte{}) {
			t.Error("returned true")
		}
	})

	t.Run("accept only file type", func(t *testing.T) {
		scenarios := []struct {
			sourceType string
			expected   bool
		}{
			{ // test file type
				sourceType: LogStreamTypeFile,
				expected:   true,
			},
			{ // test non-file format (db)
				sourceType: "db",
				expected:   false,
			},
		}

		for _, scn := range scenarios {
			if check := strategy.Accept(scn.sourceType, path, format, channels, level); check != scn.expected {
				t.Errorf("returned (%v) for the type (%s)", check, scn.sourceType)
			}
		}
	})
}

func Test_LogStreamFactoryStrategyFile_AcceptConfig(t *testing.T) {
	sourceType := LogStreamTypeFile
	path := "path"
	format := LogFormatterFormatJSON
	channels := []interface{}{"channel.1", "channel.2"}
	level := "debug"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	formatterFactory := NewLogFormatterFactory()
	strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

	t.Run("don't accept if type is missing", func(t *testing.T) {
		partial := ConfigPartial{}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if type is not a string", func(t *testing.T) {
		partial := ConfigPartial{"type": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if path is missing", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if path is not a string", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if format is missing", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if format is not a string", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path, "format": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if channels is missing", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path, "format": format}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if channels is not a list of strings", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path, "format": format, "channels": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if level is missing", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path, "format": format, "channels": channels}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if level is not a string", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path, "format": format, "channels": channels, "level": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if level is unrecognizable", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path, "format": format, "channels": channels, "level": "unknown"}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("accept the config", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path, "format": format, "channels": channels, "level": level}
		if !strategy.AcceptConfig(partial) {
			t.Error("returned false")
		}
	})
}

func Test_LogStreamFactoryStrategyFile_Create(t *testing.T) {
	path := "path"
	format := "json"
	var channels []string
	level := DEBUG

	t.Run("non-string path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		if source, err := strategy.Create(123, format, channels, level); source != nil {
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
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		if source, err := strategy.Create(path, 123, channels, level); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("non-string list channels", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		if source, err := strategy.Create(path, format, "string", level); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("non-loglevel level", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		if source, err := strategy.Create(path, format, channels, "string"); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on opening the file", func(t *testing.T) {
		expectedError := "dummy_error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)).Times(1)
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		if stream, err := strategy.Create(path, format, channels, level); stream != nil {
			_ = stream.Close()
			t.Error("returned a valid stream")
		} else if err == nil {
			t.Error("didn't returned the expected error")
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
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		if stream, err := strategy.Create(path, format, channels, level); stream != nil {
			_ = stream.Close()
			t.Error("returned a valid stream")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized format type : json" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the file stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		formatterFactory := NewLogFormatterFactory()
		_ = formatterFactory.Register(NewLogFormatterFactoryStrategyJSON())
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		if stream, err := strategy.Create(path, format, channels, level); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if stream == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch stream.(type) {
			case *LogStreamFile:
			default:
				t.Error("didn't returned a new file stream")
			}
		}
	})
}

func Test_FileStreamFactoryStrategy_CreateConfig(t *testing.T) {
	t.Run("non-string path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		conf := ConfigPartial{"path": 123}
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
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		conf := ConfigPartial{"path": "path", "format": 123}
		if source, err := strategy.CreateConfig(conf); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("non-list channels", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		conf := ConfigPartial{"path": "path", "format": "format", "channels": 123}
		if source, err := strategy.CreateConfig(conf); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("non-string level", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		conf := ConfigPartial{"path": "path", "format": "format", "channels": []interface{}{}, "level": 123}
		if source, err := strategy.CreateConfig(conf); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("non-loglevel name level", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		conf := ConfigPartial{"path": "path", "format": "format", "channels": []interface{}{}, "level": "invalid"}
		if source, err := strategy.CreateConfig(conf); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized logger level : invalid" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new level", func(t *testing.T) {
		path := "path"
		format := "json"
		channels := []interface{}{"channel1"}
		level := "debug"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		formatterFactory := NewLogFormatterFactory()
		_ = formatterFactory.Register(NewLogFormatterFactoryStrategyJSON())
		strategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)

		conf := ConfigPartial{"path": path, "format": format, "channels": channels, "level": level}
		if stream, err := strategy.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if stream == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch stream.(type) {
			case *LogStreamFile:
			default:
				t.Error("didn't returned a new file stream")
			}
		}
	})
}
