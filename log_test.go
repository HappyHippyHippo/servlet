package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

/// ---------------------------------------------------------------------------
/// LogLevel
/// ---------------------------------------------------------------------------

func Test_LogLevel(t *testing.T) {
	t.Run("levels have correct priorities", func(t *testing.T) {
		scenarios := []struct {
			lower      LogLevel
			lowerName  string
			higher     LogLevel
			higherName string
		}{
			{
				lower:      FATAL,
				lowerName:  "FATAL",
				higher:     ERROR,
				higherName: "ERROR",
			},
			{
				lower:      ERROR,
				lowerName:  "ERROR",
				higher:     WARNING,
				higherName: "WARNING",
			},
			{
				lower:      WARNING,
				lowerName:  "WARNING",
				higher:     NOTICE,
				higherName: "NOTICE",
			},
			{
				lower:      NOTICE,
				lowerName:  "NOTICE",
				higher:     INFO,
				higherName: "INFO",
			},
			{
				lower:      INFO,
				lowerName:  "INFO",
				higher:     DEBUG,
				higherName: "DEBUG",
			},
		}

		for _, scn := range scenarios {
			if scn.lower > scn.higher {
				t.Errorf("lower %s greater then %s", scn.lowerName, scn.higherName)
			}
		}
	})
}

func Test_LogLevelMap(t *testing.T) {
	t.Run("level map have correct priorities", func(t *testing.T) {
		scenarios := []struct {
			name  string
			level LogLevel
		}{
			{
				name:  "fatal",
				level: FATAL,
			},
			{
				name:  "error",
				level: ERROR,
			},
			{
				name:  "warning",
				level: WARNING,
			},
			{
				name:  "notice",
				level: NOTICE,
			},
			{
				name:  "info",
				level: INFO,
			},
			{
				name:  "debug",
				level: DEBUG,
			},
		}

		for _, scn := range scenarios {
			if scn.level != LogLevelMap[scn.name] {
				t.Errorf("(%s) did not correspond to (%v) level", scn.name, scn.level)
			}
		}
	})
}

func Test_LogLevelNameMap(t *testing.T) {
	t.Run("level map have correct priorities", func(t *testing.T) {
		scenarios := []struct {
			name  string
			level LogLevel
		}{
			{
				name:  "fatal",
				level: FATAL,
			},
			{
				name:  "error",
				level: ERROR,
			},
			{
				name:  "warning",
				level: WARNING,
			},
			{
				name:  "notice",
				level: NOTICE,
			},
			{
				name:  "info",
				level: INFO,
			},
			{
				name:  "debug",
				level: DEBUG,
			},
		}

		for _, scn := range scenarios {
			if scn.name != LogLevelNameMap[scn.level] {
				t.Errorf("(%v) did not correspond to (%s) name", scn.level, scn.name)
			}
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogFormatterFactory
/// ---------------------------------------------------------------------------

func Test_NewLogFormatterFactory(t *testing.T) {
	t.Run("new log formatter factory", func(t *testing.T) {
		if NewLogFormatterFactory() == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_LogFormatterFactory_Register(t *testing.T) {
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

		var factory *LogFormatterFactory
		_ = factory.Register(nil)
	})

	t.Run("nil strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogFormatterFactory()

		if err := factory.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'strategy' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the formatter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockLogFormatterFactoryStrategy(ctrl)
		factory := NewLogFormatterFactory()

		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.strategies[0] != strategy {
			t.Errorf("didn't stored the strategy")
		}
	})
}

func Test_LogFormatterFactory_Create(t *testing.T) {
	t.Run("unrecognized format", func(t *testing.T) {
		format := "invalid format"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogFormatterFactory()

		strategy := NewMockLogFormatterFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format).Return(false).Times(1)
		_ = factory.Register(strategy)

		if result, err := factory.Create(format); result != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized format type : invalid format" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the formatter", func(t *testing.T) {
		format := "format"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogFormatterFactory()

		formatter := NewLogJSONFormatter()
		strategy := NewMockLogFormatterFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format).Return(true).Times(1)
		strategy.EXPECT().Create().Return(formatter, nil).Times(1)
		_ = factory.Register(strategy)

		if formatter, err := factory.Create(format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(formatter, formatter) {
			t.Errorf("didn't returned the formatter")
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogJSONFormatter
/// ---------------------------------------------------------------------------

func Test_NewLogJSONFormatter(t *testing.T) {
	t.Run("new json formatter", func(t *testing.T) {
		if NewLogJSONFormatter() == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_LogJSONFormatter_Format(t *testing.T) {
	t.Run("correctly format the message", func(t *testing.T) {
		scenarios := []struct {
			level    LogLevel
			fields   map[string]interface{}
			message  string
			expected string
		}{
			{ // test level FATAL
				level:    FATAL,
				fields:   nil,
				message:  "",
				expected: `"level"\s*\:\s*"FATAL"`,
			},
			{ // test level ERROR
				level:    ERROR,
				fields:   nil,
				message:  "",
				expected: `"level"\s*\:\s*"ERROR"`,
			},
			{ // test level WARNING
				level:    WARNING,
				fields:   nil,
				message:  "",
				expected: `"level"\s*\:\s*"WARNING"`,
			},
			{ // test level NOTICE
				level:    NOTICE,
				fields:   nil,
				message:  "",
				expected: `"level"\s*\:\s*"NOTICE"`,
			},
			{ // test level INFO
				level:    INFO,
				fields:   nil,
				message:  "",
				expected: `"level"\s*\:\s*"INFO"`,
			},
			{ // test level DEBUG
				level:    DEBUG,
				fields:   nil,
				message:  "",
				expected: `"level"\s*\:\s*"DEBUG"`,
			},
			{ // test fields (single value)
				level:    DEBUG,
				fields:   map[string]interface{}{"field1": "value1"},
				message:  "",
				expected: `"field1"\s*\:\s*"value1"`,
			},
			{ // test fields (multiple value)
				level:    DEBUG,
				fields:   map[string]interface{}{"field1": "value1", "field2": "value2"},
				message:  "",
				expected: `"field1"\s*\:\s*"value1"|"field2"\s*\:\s*"value2"`,
			},
			{ // test message
				level:    DEBUG,
				fields:   nil,
				message:  "My_message",
				expected: `"message"\s*\:\s*"My_message"`,
			},
		}

		formatter := NewLogJSONFormatter()

		for _, scn := range scenarios {
			result := formatter.Format(scn.level, scn.message, scn.fields)
			matched, _ := regexp.Match(scn.expected, []byte(result))
			if !matched {
				t.Errorf("didn't validated (%s) output", result)
			}
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogJSONFormatterFactoryStrategy
/// ---------------------------------------------------------------------------

func Test_NewLogJSONFormatterFactoryStrategy(t *testing.T) {
	t.Run("new json formatter factory strategy", func(t *testing.T) {
		if NewLogJSONFormatterFactoryStrategy() == nil {
			t.Errorf("didn't returned a valid reference")
		}
	})
}

func Test_LogJSONFormatterFactoryStrategy_Accept(t *testing.T) {
	t.Run("accept only json format", func(t *testing.T) {
		scenarios := []struct {
			format   string
			expected bool
		}{
			{ // test json format
				format:   LogFormatterFormatJSON,
				expected: true,
			},
			{ // test non-json format (yaml)
				format:   "yaml",
				expected: false,
			},
		}

		strategy := NewLogJSONFormatterFactoryStrategy()

		for _, scn := range scenarios {
			if check := strategy.Accept(scn.format); check != scn.expected {
				t.Errorf("returned (%v) for the (%s) format", check, scn.format)
			}
		}
	})
}

func Test_LogJSONFormatterFactoryStrategy_Create(t *testing.T) {
	strategy := NewLogJSONFormatterFactoryStrategy()

	t.Run("create json formatter", func(t *testing.T) {
		if formatter, err := strategy.Create(); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if formatter == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch formatter.(type) {
			case *LogJSONFormatter:
			default:
				t.Errorf("didn't returned a new json formatter")
			}
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogBaseStream
/// ---------------------------------------------------------------------------

func Test_LogBaseStream_HasChannel(t *testing.T) {
	channels := []string{"channel.1", "channel.2"}
	stream := &LogBaseStream{nil, channels, WARNING}

	t.Run("check the channel registration", func(t *testing.T) {

		if !stream.HasChannel("channel.1") {
			t.Error("'channel.1' channel was not found")
		} else if !stream.HasChannel("channel.2") {
			t.Error("'channel.2' channel was not found")
		} else if stream.HasChannel("channel.3") {
			t.Error("'channel.3' channel was found")
		}
	})
}

func Test_LogBaseStream_ListChannels(t *testing.T) {
	channels := []string{"channel.1", "channel.2"}
	stream := &LogBaseStream{nil, channels, WARNING}

	t.Run("list the registered channels", func(t *testing.T) {
		if result := stream.ListChannels(); !reflect.DeepEqual(result, channels) {
			t.Errorf("returned the (%v) list of channels", result)
		}
	})
}

func Test_LogBaseStream_AddChannel(t *testing.T) {
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

		var stream *LogBaseStream
		stream.AddChannel("channel")
	})

	t.Run("register a new channel", func(t *testing.T) {
		scenarios := []struct {
			state struct {
				channels []string
				level    LogLevel
			}
			channel  string
			expected []string
		}{
			{ // adding into a empty list
				state: struct {
					channels []string
					level    LogLevel
				}{
					channels: []string{},
					level:    DEBUG,
				},
				channel:  "channel.1",
				expected: []string{"channel.1"},
			},
			{ // adding should keep sorting
				state: struct {
					channels []string
					level    LogLevel
				}{
					channels: []string{"channel.1", "channel.3"},
					level:    DEBUG,
				},
				channel:  "channel.2",
				expected: []string{"channel.1", "channel.2", "channel.3"},
			},
			{ // adding an already existent should result in a no-op
				state: struct {
					channels []string
					level    LogLevel
				}{
					channels: []string{"channel.1", "channel.2", "channel.3"},
					level:    DEBUG,
				},
				channel:  "channel.2",
				expected: []string{"channel.1", "channel.2", "channel.3"},
			},
		}

		for _, scn := range scenarios {
			stream := &LogBaseStream{nil, scn.state.channels, scn.state.level}
			stream.AddChannel(scn.channel)

			if result := stream.ListChannels(); !reflect.DeepEqual(result, scn.expected) {
				t.Errorf("returned the (%v) list of channels", result)
			}
		}
	})
}

func Test_LogBaseStream_RemoveChannel(t *testing.T) {
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

		var stream *LogBaseStream
		stream.RemoveChannel("channel")
	})

	t.Run("unregister a channel", func(t *testing.T) {
		scenarios := []struct {
			state struct {
				channels []string
				level    LogLevel
			}
			channel  string
			expected []string
		}{
			{ // removing from an empty list
				state: struct {
					channels []string
					level    LogLevel
				}{
					channels: []string{},
					level:    DEBUG,
				},
				channel:  "channel.1",
				expected: []string{},
			},
			{ // removing a non existing channel
				state: struct {
					channels []string
					level    LogLevel
				}{
					channels: []string{"channel.1", "channel.3"},
					level:    DEBUG,
				},
				channel:  "channel.2",
				expected: []string{"channel.1", "channel.3"},
			},
			{ // removing an existing channel
				state: struct {
					channels []string
					level    LogLevel
				}{
					channels: []string{"channel.1", "channel.2", "channel.3"},
					level:    DEBUG,
				},
				channel:  "channel.2",
				expected: []string{"channel.1", "channel.3"},
			},
		}

		for _, scn := range scenarios {
			stream := &LogBaseStream{nil, scn.state.channels, scn.state.level}
			stream.RemoveChannel(scn.channel)

			if result := stream.ListChannels(); !reflect.DeepEqual(result, scn.expected) {
				t.Errorf("returned the (%v) list of channels", result)
			}
		}
	})
}

func Test_LogBaseStream_Level(t *testing.T) {
	level := WARNING
	stream := &LogBaseStream{nil, []string{}, level}

	t.Run("retrieve the filtering level", func(t *testing.T) {
		if result := stream.Level(); result != level {
			t.Errorf("returned the (%v) level", result)
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogStreamFactory
/// ---------------------------------------------------------------------------

func Test_NewLogStreamFactory(t *testing.T) {
	t.Run("new config stream factory", func(t *testing.T) {
		if NewLogStreamFactory() == nil {
			t.Errorf("didn't returned a valid reference")
		}
	})
}

func Test_LogStreamFactory_Register(t *testing.T) {
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

		var factory *LogStreamFactory
		_ = factory.Register(nil)
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	strategy := NewMockLogStreamFactoryStrategy(ctrl)
	factory := NewLogStreamFactory()

	t.Run("nil strategy", func(t *testing.T) {
		if err := factory.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if check := err.Error(); check != "invalid nil 'strategy' argument" {
			t.Errorf("return the (%v) error", check)
		}
	})

	t.Run("register the stream factory strategy", func(t *testing.T) {
		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.strategies[0] != strategy {
			t.Error("didn't stored the strategy")
		}
	})
}

func Test_LogStreamFactory_Create(t *testing.T) {
	sourceType := "type"
	path := "path"
	format := "format"

	t.Run("unrecognized format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogStreamFactory()

		strategy := NewMockLogStreamFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(sourceType, path, format).Return(false).Times(1)
		_ = factory.Register(strategy)

		if stream, err := factory.Create(sourceType, path, format); stream != nil {
			t.Error("returned an valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized stream type : type" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the config stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogStreamFactory()

		stream := NewMockLogStream(ctrl)
		strategy := NewMockLogStreamFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(sourceType, path, format).Return(true).Times(1)
		strategy.EXPECT().Create(path, format).Return(stream, nil).Times(1)
		_ = factory.Register(strategy)

		if stream, err := factory.Create(sourceType, path, format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(stream, stream) {
			t.Error("didn't returned the created stream")
		}
	})
}

func Test_LogStreamFactory_CreateConfig(t *testing.T) {
	t.Run("unrecognized type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogStreamFactory()

		conf := ConfigPartial{}
		strategy := NewMockLogStreamFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(false).Times(1)
		_ = factory.Register(strategy)

		if stream, err := factory.CreateConfig(conf); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != fmt.Sprintf("unrecognized stream config : %v", conf) {
			t.Errorf("returned the (%v) error", err)
		} else if stream != nil {
			t.Error("returned a config stream")
		}
	})

	t.Run("create the config stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogStreamFactory()

		conf := ConfigPartial{}
		stream := NewMockLogStream(ctrl)
		strategy := NewMockLogStreamFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(true).Times(1)
		strategy.EXPECT().CreateConfig(conf).Return(stream, nil).Times(1)
		_ = factory.Register(strategy)

		if stream, err := factory.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(stream, stream) {
			t.Error("didn't returned the created stream")
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogFileStream
/// ---------------------------------------------------------------------------

func Test_NewLogFileStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := NewMockWriter(ctrl)
	writer.EXPECT().Close().Times(1)
	formatter := NewMockLogFormatter(ctrl)
	var channels []string
	level := WARNING

	t.Run("nil writer", func(t *testing.T) {
		if stream, err := NewLogFileStream(nil, formatter, channels, level); stream != nil {
			_ = stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'writer' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil formatter", func(t *testing.T) {
		if stream, err := NewLogFileStream(writer, nil, channels, level); stream != nil {
			_ = stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'formatter' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new file stream", func(t *testing.T) {
		if stream, err := NewLogFileStream(writer, formatter, []string{}, WARNING); stream == nil {
			t.Error("didn't returned a valid reference")
		} else {
			_ = stream.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}

func Test_LogFileStream_Close(t *testing.T) {
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

		var stream *LogFileStream
		_ = stream.Close()
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := NewMockWriter(ctrl)
	writer.EXPECT().Close().Times(1)
	formatter := NewMockLogFormatter(ctrl)

	stream, _ := NewLogFileStream(writer, formatter, []string{}, WARNING)

	t.Run("call the close on the writer only once", func(t *testing.T) {
		_ = stream.Close()
		_ = stream.Close()
	})
}

func Test_LogFileStream_Signal(t *testing.T) {
	t.Run("signal message to the writer", func(t *testing.T) {
		scenarios := []struct {
			state struct {
				channels []string
				level    LogLevel
			}
			call struct {
				level   LogLevel
				channel string
				fields  map[string]interface{}
				message string
			}
			callTimes int
			expected  string
		}{
			{ // signal through a valid channel with a not filtered level
				state: struct {
					channels []string
					level    LogLevel
				}{
					channels: []string{"dummy_channel"},
					level:    WARNING,
				},
				call: struct {
					level   LogLevel
					channel string
					fields  map[string]interface{}
					message string
				}{
					level:   FATAL,
					channel: "dummy_channel",
					fields:  map[string]interface{}{},
					message: "dummy_message",
				},
				callTimes: 1,
				expected:  `{"message" : "dummy_message"}`,
			},
			{ // signal through a valid channel with a filtered level
				state: struct {
					channels []string
					level    LogLevel
				}{
					channels: []string{"dummy_channel"},
					level:    WARNING,
				},
				call: struct {
					level   LogLevel
					channel string
					fields  map[string]interface{}
					message string
				}{
					level:   DEBUG,
					channel: "dummy_channel",
					fields:  map[string]interface{}{},
					message: "dummy_message",
				},
				callTimes: 0,
				expected:  `{"message" : "dummy_message"}`,
			},
			{ // signal through a valid channel with a unregistered channel
				state: struct {
					channels []string
					level    LogLevel
				}{
					channels: []string{"dummy_channel"},
					level:    WARNING,
				},
				call: struct {
					level   LogLevel
					channel string
					fields  map[string]interface{}
					message string
				}{
					level:   FATAL,
					channel: "not_a_valid_dummy_channel",
					fields:  map[string]interface{}{},
					message: "dummy_message",
				},
				callTimes: 0,
				expected:  `{"message" : "dummy_message"}`,
			},
		}

		for _, scn := range scenarios {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			writer := NewMockWriter(ctrl)
			writer.EXPECT().Close().Times(1)
			writer.EXPECT().Write([]byte(scn.expected + "\n")).Times(scn.callTimes)
			formatter := NewMockLogFormatter(ctrl)
			formatter.EXPECT().Format(scn.call.level, scn.call.message, scn.call.fields).Return(scn.expected).Times(scn.callTimes)
			stream, _ := NewLogFileStream(writer, formatter, scn.state.channels, scn.state.level)
			defer func() { _ = stream.Close() }()

			if err := stream.Signal(scn.call.channel, scn.call.level, scn.call.message, scn.call.fields); err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}

func Test_LogFileStream_Broadcast(t *testing.T) {
	t.Run("broadcast message to the writer", func(t *testing.T) {
		scenarios := []struct {
			state struct {
				channels []string
				level    LogLevel
			}
			call struct {
				level   LogLevel
				fields  map[string]interface{}
				message string
			}
			callTimes int
			expected  string
		}{
			{ // broadcast through a valid channel with a not filtered level
				state: struct {
					channels []string
					level    LogLevel
				}{
					channels: []string{"dummy_channel"},
					level:    WARNING,
				},
				call: struct {
					level   LogLevel
					fields  map[string]interface{}
					message string
				}{
					level:   FATAL,
					fields:  map[string]interface{}{},
					message: "dummy_message",
				},
				callTimes: 1,
				expected:  `{"message" : "dummy_message"}`,
			},
			{ // broadcast through a valid channel with a filtered level
				state: struct {
					channels []string
					level    LogLevel
				}{
					channels: []string{"dummy_channel"},
					level:    WARNING,
				},
				call: struct {
					level   LogLevel
					fields  map[string]interface{}
					message string
				}{
					level:   DEBUG,
					fields:  map[string]interface{}{},
					message: "dummy_message",
				},
				callTimes: 0,
				expected:  `{"message" : "dummy_message"}`,
			},
		}

		for _, scn := range scenarios {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			writer := NewMockWriter(ctrl)
			writer.EXPECT().Close().Times(1)
			writer.EXPECT().Write([]byte(scn.expected + "\n")).Times(scn.callTimes)
			formatter := NewMockLogFormatter(ctrl)
			formatter.EXPECT().Format(scn.call.level, scn.call.message, scn.call.fields).Return(scn.expected).Times(scn.callTimes)
			stream, _ := NewLogFileStream(writer, formatter, scn.state.channels, scn.state.level)
			defer func() { _ = stream.Close() }()

			if err := stream.Broadcast(scn.call.level, scn.call.message, scn.call.fields); err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogFileStreamFactoryStrategy
/// ---------------------------------------------------------------------------

func Test_NewLogFileStreamFactoryStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	formatterFactory := NewLogFormatterFactory()

	t.Run("nil file system adapter", func(t *testing.T) {
		if strategy, err := NewLogFileStreamFactoryStrategy(nil, formatterFactory); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil formatter factory", func(t *testing.T) {
		if strategy, err := NewLogFileStreamFactoryStrategy(fileSystem, nil); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'formatterFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new file stream factory strategy", func(t *testing.T) {
		if strategy, err := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory); strategy == nil {
			t.Errorf("didn't returned a valid reference")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_LogFileStreamFactoryStrategy_Accept(t *testing.T) {
	path := "path"
	format := "format"
	var channels []string
	level := DEBUG

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	formatterFactory := NewLogFormatterFactory()
	strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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

func Test_LogFileStreamFactoryStrategy_AcceptConfig(t *testing.T) {
	sourceType := LogStreamTypeFile
	path := "path"
	format := LogFormatterFormatJSON
	channels := []interface{}{"channel.1", "channel.2"}
	level := "debug"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	formatterFactory := NewLogFormatterFactory()
	strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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

func Test_LogFileStreamFactoryStrategy_Create(t *testing.T) {
	path := "path"
	format := "json"
	var channels []string
	level := DEBUG

	t.Run("non-string path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewLogFormatterFactory()
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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
		_ = formatterFactory.Register(NewLogJSONFormatterFactoryStrategy())
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

		if stream, err := strategy.Create(path, format, channels, level); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if stream == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch stream.(type) {
			case *LogFileStream:
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
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

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
		_ = formatterFactory.Register(NewLogJSONFormatterFactoryStrategy())
		strategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)

		conf := ConfigPartial{"path": path, "format": format, "channels": channels, "level": level}
		if stream, err := strategy.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if stream == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch stream.(type) {
			case *LogFileStream:
			default:
				t.Error("didn't returned a new file stream")
			}
		}
	})
}

/// ---------------------------------------------------------------------------
/// Log
/// ---------------------------------------------------------------------------

func Test_NewLog(t *testing.T) {
	t.Run("new logger", func(t *testing.T) {
		if logger := NewLog(); logger == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_Log_Close(t *testing.T) {
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

		var log *Log
		_ = log.Close()
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := NewLog()

	id1 := "stream.1"
	stream1 := NewMockLogStream(ctrl)
	stream1.EXPECT().Close().Times(1)
	_ = logger.AddStream(id1, stream1)

	id2 := "stream.2"
	stream2 := NewMockLogStream(ctrl)
	stream2.EXPECT().Close().Times(1)
	_ = logger.AddStream(id2, stream2)

	t.Run("execute close process", func(t *testing.T) {
		_ = logger.Close()

		if logger.HasStream(id1) {
			t.Error("didn't removed the stream")
		}
		if logger.HasStream(id2) {
			t.Error("didn't removed the stream")
		}
	})
}

func Test_Log_Signal(t *testing.T) {
	id1 := "stream.1"
	id2 := "stream.2"
	channel := "channel"
	level := WARNING
	fields := map[string]interface{}{"field": "value"}
	message := "message"

	t.Run("propagate to all streams", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Signal(channel, level, message, fields).Return(nil).Times(1)
		stream1.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id1, stream1)

		stream2 := NewMockLogStream(ctrl)
		stream2.EXPECT().Signal(channel, level, message, fields).Return(nil).Times(1)
		stream2.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id2, stream2)

		if err := logger.Signal(channel, level, message, fields); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("return on the first error", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Signal(channel, level, message, fields).Return(fmt.Errorf(expectedError)).AnyTimes()
		stream1.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id1, stream1)

		stream2 := NewMockLogStream(ctrl)
		stream2.EXPECT().Signal(channel, level, message, fields).Return(nil).AnyTimes()
		stream2.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id2, stream2)

		if err := logger.Signal(channel, level, message, fields); err == nil {
			t.Error("didn't returned the expected  error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_Log_Broadcast(t *testing.T) {
	id1 := "stream.1"
	id2 := "stream.2"
	level := WARNING
	fields := map[string]interface{}{"field": "value"}
	message := "message"

	t.Run("propagate to all streams", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Broadcast(level, message, fields).Return(nil).Times(1)
		stream1.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id1, stream1)

		stream2 := NewMockLogStream(ctrl)
		stream2.EXPECT().Broadcast(level, message, fields).Return(nil).Times(1)
		stream2.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id2, stream2)

		if err := logger.Broadcast(level, message, fields); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("return on the first error", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Broadcast(level, message, fields).Return(fmt.Errorf(expectedError)).AnyTimes()
		stream1.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id1, stream1)

		stream2 := NewMockLogStream(ctrl)
		stream2.EXPECT().Broadcast(level, message, fields).Return(nil).AnyTimes()
		stream2.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id2, stream2)

		if err := logger.Broadcast(level, message, fields); err == nil {
			t.Error("didn't returned the expected  error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_Log_HasStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := NewLog()
	defer func() { _ = logger.Close() }()

	id1 := "stream.1"
	stream1 := NewMockLogStream(ctrl)
	stream1.EXPECT().Close().Return(nil).Times(1)
	_ = logger.AddStream(id1, stream1)

	id2 := "stream.2"
	stream2 := NewMockLogStream(ctrl)
	stream2.EXPECT().Close().Return(nil).Times(1)
	_ = logger.AddStream(id2, stream2)

	id3 := "stream.3"

	t.Run("check the registration of a stream", func(t *testing.T) {
		if !logger.HasStream(id1) {
			t.Errorf("returned false")
		}
		if !logger.HasStream(id2) {
			t.Errorf("returned false")
		}
		if logger.HasStream(id3) {
			t.Errorf("returned true")
		}
	})
}

func Test_Log_AddStream(t *testing.T) {
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

		var log *Log
		_ = log.AddStream("id", nil)
	})

	t.Run("error if nil stream", func(t *testing.T) {
		logger := NewLog()
		defer func() { _ = logger.Close() }()

		if err := logger.AddStream("id", nil); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'stream' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error if id is duplicate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		id := "stream"
		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Close().Return(nil).Times(1)

		stream2 := NewMockLogStream(ctrl)
		_ = logger.AddStream(id, stream1)

		if err := logger.AddStream(id, stream2); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "duplicate id : stream" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register a new stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		id := "stream"
		stream := NewMockLogStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)

		if err := logger.AddStream(id, stream); err != nil {
			t.Errorf("resulted the (%v) error", err)
		} else if check := logger.Stream(id); !reflect.DeepEqual(check, stream) {
			t.Errorf("didn't stored the stream")
		}
	})
}

func Test_Log_RemoveStream(t *testing.T) {
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

		var log *Log
		log.RemoveStream("id")
	})

	t.Run("unregister a stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		id := "stream"
		stream := NewMockLogStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)

		_ = logger.AddStream(id, stream)
		logger.RemoveStream(id)

		if logger.HasStream(id) {
			t.Errorf("dnd't removed the stream")
		}
	})
}

func Test_Log_Stream(t *testing.T) {
	t.Run("nil on a non-existing stream", func(t *testing.T) {
		logger := NewLog()
		defer func() { _ = logger.Close() }()

		if result := logger.Stream("invalid id"); result != nil {
			t.Errorf("returned a valid stream")
		}
	})

	t.Run("retrieve the requested stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		id := "stream"
		stream := NewMockLogStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id, stream)

		if check := logger.Stream(id); !reflect.DeepEqual(check, stream) {
			t.Errorf("didn0t retrieved the stored stream")
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogLoader
/// ---------------------------------------------------------------------------

func Test_NewLogLoader(t *testing.T) {
	logger := NewLog()
	streamFactory := NewLogStreamFactory()

	t.Run("error when missing the logger", func(t *testing.T) {
		if loader, err := NewLogLoader(nil, streamFactory); loader != nil {
			t.Errorf("return a valid reference")
		} else if err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'logger' argument" {
			t.Errorf("returned the (%v)) error", err)
		}
	})

	t.Run("error when missing the logger stream factory", func(t *testing.T) {
		if loader, err := NewLogLoader(logger, nil); loader != nil {
			t.Errorf("return a valid reference")
		} else if err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'streamFactory' argument" {
			t.Errorf("returned the (%v)) error", err)
		}
	})

	t.Run("create loader", func(t *testing.T) {
		if loader, err := NewLogLoader(logger, streamFactory); loader == nil {
			t.Errorf("didn't returned a valid reference")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_LogLoader_Load(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		if err := loader.Load(nil); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'config' argument" {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("no-op if stream list is missing", func(t *testing.T) {
		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		config, _ := NewConfig(0 * time.Second)

		if err := loader.Load(config); err != nil {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("no-op if stream list is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		conf := ConfigPartial{"log": ConfigPartial{"sources": []interface{}{}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)

		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("source", 0, source)

		if err := loader.Load(config); err != nil {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("error if stream list is not a list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		conf := ConfigPartial{"log": ConfigPartial{"streams": 123}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)

		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("source", 0, source)

		if err := loader.Load(config); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("missing stream id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		streamConfig := ConfigPartial{}
		conf := ConfigPartial{"log": ConfigPartial{"streams": []interface{}{streamConfig}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)

		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("source", 0, source)

		if err := loader.Load(config); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("invalid stream id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		streamConfig := ConfigPartial{"id": 123}
		conf := ConfigPartial{"log": ConfigPartial{"streams": []interface{}{streamConfig}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)
		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("source", 0, source)

		if err := loader.Load(config); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("error creating stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		streamConfig := ConfigPartial{"id": "id"}
		conf := ConfigPartial{"log": ConfigPartial{"streams": []interface{}{streamConfig}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)
		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("source", 0, source)

		if err := loader.Load(config); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if strings.Index(err.Error(), "unrecognized stream config :") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("error storing stream", func(t *testing.T) {
		streamConfig := ConfigPartial{
			"id":       "id",
			"type":     "file",
			"path":     "path",
			"format":   "json",
			"channels": []interface{}{},
			"level":    "debug"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile("path", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		formatterFactory := NewLogFormatterFactory()
		_ = formatterFactory.Register(NewLogJSONFormatterFactoryStrategy())

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		fileStreamFactoryStrategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)
		_ = streamFactory.Register(fileStreamFactoryStrategy)
		loader, _ := NewLogLoader(logger, streamFactory)

		conf := ConfigPartial{"log": ConfigPartial{"streams": []interface{}{streamConfig}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)

		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("id", 0, source)

		writer := NewMockWriter(ctrl)
		formatter := NewLogJSONFormatter()
		fileLogger, _ := NewLogFileStream(writer, formatter, []string{}, FATAL)
		_ = logger.AddStream("id", fileLogger)

		if err := loader.Load(config); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if strings.Index(err.Error(), "duplicate id :") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("register stream", func(t *testing.T) {
		streamConfig := ConfigPartial{
			"id":       "id",
			"type":     "file",
			"path":     "path",
			"format":   "json",
			"channels": []interface{}{},
			"level":    "debug"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile("path", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		formatterFactory := NewLogFormatterFactory()
		_ = formatterFactory.Register(NewLogJSONFormatterFactoryStrategy())

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		fileStreamFactoryStrategy, _ := NewLogFileStreamFactoryStrategy(fileSystem, formatterFactory)
		_ = streamFactory.Register(fileStreamFactoryStrategy)
		loader, _ := NewLogLoader(logger, streamFactory)

		conf := ConfigPartial{"log": ConfigPartial{"streams": []interface{}{streamConfig}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)

		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("id", 0, source)

		if err := loader.Load(config); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !logger.HasStream("id") {
			t.Error("didn't stored the loaded stream")
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogParams
/// ---------------------------------------------------------------------------

func Test_NewLogParams(t *testing.T) {
	t.Run("new parameters", func(t *testing.T) {
		parameters := NewLogParams()

		if value := parameters.LoggerID; value != ContainerLoggerID {
			t.Errorf("stored (%v) logger ID", value)
		} else if value := parameters.FileSystemID; value != ContainerFileSystemID {
			t.Errorf("stored (%v) file sytem ID", value)
		} else if value := parameters.ConfigID; value != ContainerConfigID {
			t.Errorf("stored (%v) config ID", value)
		} else if value := parameters.FormatterFactoryID; value != ContainerLogFormatterFactoryID {
			t.Errorf("stored (%v) formatter factory ID", value)
		} else if value := parameters.StreamFactoryID; value != ContainerLogStreamFactoryID {
			t.Errorf("stored (%v) stream factory ID", value)
		} else if value := parameters.LoaderID; value != ContainerLogLoaderID {
			t.Errorf("stored (%v) loader ID", value)
		}
	})

	t.Run("with the env logger ID", func(t *testing.T) {
		loggerID := "logger_id"
		_ = os.Setenv(EnvContainerLoggerID, loggerID)
		defer func() { _ = os.Setenv(EnvContainerLoggerID, "") }()

		parameters := NewLogParams()
		if value := parameters.LoggerID; value != loggerID {
			t.Errorf("stored (%v) logger ID", value)
		}
	})

	t.Run("with the env file system ID", func(t *testing.T) {
		fileSystemID := "file_system_id"
		_ = os.Setenv(EnvContainerFileSystemID, fileSystemID)
		defer func() { _ = os.Setenv(EnvContainerFileSystemID, "") }()

		parameters := NewLogParams()
		if value := parameters.FileSystemID; value != fileSystemID {
			t.Errorf("stored (%v) file system ID", value)
		}
	})

	t.Run("with the env config ID", func(t *testing.T) {
		configID := "config_id"
		_ = os.Setenv(EnvContainerConfigID, configID)
		defer func() { _ = os.Setenv(EnvContainerConfigID, "") }()

		parameters := NewLogParams()
		if value := parameters.ConfigID; value != configID {
			t.Errorf("stored (%v) config ID", value)
		}
	})

	t.Run("with the env formatter factory ID", func(t *testing.T) {
		formatterFactoryID := "formatter_factory_id"
		_ = os.Setenv(EnvContainerLogFormatterFactoryID, formatterFactoryID)
		defer func() { _ = os.Setenv(EnvContainerLogFormatterFactoryID, "") }()

		parameters := NewLogParams()
		if value := parameters.FormatterFactoryID; value != formatterFactoryID {
			t.Errorf("stored (%v) formatter factory ID", value)
		}
	})

	t.Run("with the env stream factory ID", func(t *testing.T) {
		streamFactoryID := "stream_factory_id"
		_ = os.Setenv(EnvContainerLogStreamFactoryID, streamFactoryID)
		defer func() { _ = os.Setenv(EnvContainerLogStreamFactoryID, "") }()

		parameters := NewLogParams()
		if value := parameters.StreamFactoryID; value != streamFactoryID {
			t.Errorf("stored (%v) stream factory ID", value)
		}
	})

	t.Run("with the env loader ID", func(t *testing.T) {
		loaderID := "loader_id"
		_ = os.Setenv(EnvContainerLogLoaderID, loaderID)
		defer func() { _ = os.Setenv(EnvContainerLogLoaderID, "") }()

		parameters := NewLogParams()
		if value := parameters.LoaderID; value != loaderID {
			t.Errorf("stored (%v) loader ID", value)
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogProvider
/// ---------------------------------------------------------------------------

func Test_NewLogProvider(t *testing.T) {
	t.Run("without params", func(t *testing.T) {
		if provider := NewLogProvider(nil); provider == nil {
			t.Error("didn't returned a valid reference")
		} else if !reflect.DeepEqual(NewLogParams(), provider.params) {
			t.Errorf("stored the (%v) parameters", provider.params)
		}
	})

	t.Run("with defined params", func(t *testing.T) {
		params := NewLogParams()
		if provider := NewLogProvider(params); provider == nil {
			t.Error("didn't returned a valid reference")
		} else if params != provider.params {
			t.Errorf("stored the (%v) parameters", provider.params)
		}
	})
}

func Test_LogProvider_Register(t *testing.T) {
	t.Run("nil container", func(t *testing.T) {
		provider := NewLogProvider(nil)
		if err := provider.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'container' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register components", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewLogProvider(nil)

		if err := provider.Register(container); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !container.Has(ContainerLogFormatterFactoryID) {
			t.Error("didnt registered the log formatter factory", err)
		} else if !container.Has(ContainerLogStreamFactoryID) {
			t.Error("didnt registered the log stream factory", err)
		} else if !container.Has(ContainerLoggerID) {
			t.Error("didnt registered the logger", err)
		} else if !container.Has(ContainerLogLoaderID) {
			t.Error("didnt registered the log loader", err)
		}
	})

	t.Run("retrieving log formatter factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		if formatterFactory, err := container.Get(ContainerLogFormatterFactoryID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if formatterFactory == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch formatterFactory.(type) {
			case *LogFormatterFactory:
			default:
				t.Error("didn't returned a formatter factory reference")
			}
		}
	})

	t.Run("error retrieving file system", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if factory, err := container.Get(ContainerLogStreamFactoryID); factory != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid file system", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if factory, err := container.Get(ContainerLogStreamFactoryID); factory != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving decoder factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLogFormatterFactoryID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if factory, err := container.Get(ContainerLogStreamFactoryID); factory != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid decoder factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLogFormatterFactoryID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if factory, err := container.Get(ContainerLogStreamFactoryID); factory != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("retrieving log stream factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		if streamFactory, err := container.Get(ContainerLogStreamFactoryID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if streamFactory == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch streamFactory.(type) {
			case *LogStreamFactory:
			default:
				t.Error("didn't returned a stream factory reference")
			}
		}
	})

	t.Run("retrieving logger", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		if logger, err := container.Get(ContainerLoggerID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if logger == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch logger.(type) {
			case *Log:
			default:
				t.Error("didn't returned a logger reference")
			}
		}
	})

	t.Run("error retrieving logger", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLoggerID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if loader, err := container.Get(ContainerLogLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLoggerID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if loader, err := container.Get(ContainerLogLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving stream factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLogStreamFactoryID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if loader, err := container.Get(ContainerLogLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid source factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		_ = container.Add(ContainerLogStreamFactoryID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if loader, err := container.Get(ContainerLogLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("retrieving log loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)
		_ = NewLogProvider(nil).Register(container)

		if loader, err := container.Get(ContainerLogLoaderID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if loader == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch loader.(type) {
			case *LogLoader:
			default:
				t.Error("didn't returned a loader reference")
			}
		}
	})
}

func Test_LogProvider_Boot(t *testing.T) {
	t.Run("error retrieving loader", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewLogProvider(nil)
		_ = provider.Register(container)
		_ = container.Add(ContainerLogLoaderID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving config", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		provider := NewLogProvider(nil)
		_ = provider.Register(container)
		_ = container.Add(ContainerConfigID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		provider := NewLogProvider(nil)
		_ = provider.Register(container)
		_ = container.Add(ContainerLogLoaderID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		provider := NewLogProvider(nil)
		_ = provider.Register(container)
		_ = container.Add(ContainerConfigID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		provider := NewLogProvider(nil)
		_ = provider.Register(container)

		if err := provider.Boot(container); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}
