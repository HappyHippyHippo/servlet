package servlet

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_NewLogStreamFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := NewMockWriter(ctrl)
	writer.EXPECT().Close().Times(1)
	formatter := NewMockLogFormatter(ctrl)
	var channels []string
	level := WARNING

	t.Run("nil writer", func(t *testing.T) {
		if stream, err := NewLogStreamFile(nil, formatter, channels, level); stream != nil {
			_ = stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'writer' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil formatter", func(t *testing.T) {
		if stream, err := NewLogStreamFile(writer, nil, channels, level); stream != nil {
			_ = stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'formatter' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new file stream", func(t *testing.T) {
		if stream, err := NewLogStreamFile(writer, formatter, []string{}, WARNING); stream == nil {
			t.Error("didn't returned a valid reference")
		} else {
			_ = stream.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}

func Test_LogStreamFile_Close(t *testing.T) {
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

		var stream *LogStreamFile
		_ = stream.Close()
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := NewMockWriter(ctrl)
	writer.EXPECT().Close().Times(1)
	formatter := NewMockLogFormatter(ctrl)

	stream, _ := NewLogStreamFile(writer, formatter, []string{}, WARNING)

	t.Run("call the close on the writer only once", func(t *testing.T) {
		_ = stream.Close()
		_ = stream.Close()
	})
}

func Test_LogStreamFile_Signal(t *testing.T) {
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
			stream, _ := NewLogStreamFile(writer, formatter, scn.state.channels, scn.state.level)
			defer func() { _ = stream.Close() }()

			if err := stream.Signal(scn.call.channel, scn.call.level, scn.call.message, scn.call.fields); err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}

func Test_LogStreamFile_Broadcast(t *testing.T) {
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
			stream, _ := NewLogStreamFile(writer, formatter, scn.state.channels, scn.state.level)
			defer func() { _ = stream.Close() }()

			if err := stream.Broadcast(scn.call.level, scn.call.message, scn.call.fields); err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}
