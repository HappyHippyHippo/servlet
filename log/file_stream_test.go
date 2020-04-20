package log

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewFileStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := NewMockWriter(ctrl)
	writer.EXPECT().Close().Times(1)
	formatter := NewMockFormatter(ctrl)
	channels := []string{}
	level := WARNING

	t.Run("error when missing writer", func(t *testing.T) {
		if stream, err := NewFileStream(nil, formatter, channels, level); stream != nil {
			stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'writer' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error when missing formatter", func(t *testing.T) {
		if stream, err := NewFileStream(writer, nil, channels, level); stream != nil {
			stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'formatter' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("creates a new file stream", func(t *testing.T) {
		if stream, err := NewFileStream(writer, formatter, []string{}, WARNING); stream == nil {
			t.Errorf("didn't return a valid reference")
		} else {
			stream.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}

func Test_FileStream_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := NewMockWriter(ctrl)
	writer.EXPECT().Close().Times(1)
	formatter := NewMockFormatter(ctrl)

	stream, _ := NewFileStream(writer, formatter, []string{}, WARNING)

	t.Run("call the close on the writer only once", func(t *testing.T) {
		stream.Close()
		stream.Close()
	})
}

func Test_FileStream_Signal(t *testing.T) {
	t.Run("signal message to the writer", func(t *testing.T) {
		scenarios := []struct {
			state struct {
				channels []string
				level    Level
			}
			call struct {
				level   Level
				channel string
				fields  F
				message string
			}
			callTimes int
			expected  string
		}{
			{ // signal through a valid channel with a not filtered level
				state: struct {
					channels []string
					level    Level
				}{
					channels: []string{"dummy_channel"},
					level:    WARNING,
				},
				call: struct {
					level   Level
					channel string
					fields  F
					message string
				}{
					level:   FATAL,
					channel: "dummy_channel",
					fields:  F{},
					message: "dummy_message",
				},
				callTimes: 1,
				expected:  `{"message" : "dummy_message"}`,
			},
			{ // signal through a valid channel with a filtered level
				state: struct {
					channels []string
					level    Level
				}{
					channels: []string{"dummy_channel"},
					level:    WARNING,
				},
				call: struct {
					level   Level
					channel string
					fields  F
					message string
				}{
					level:   DEBUG,
					channel: "dummy_channel",
					fields:  F{},
					message: "dummy_message",
				},
				callTimes: 0,
				expected:  `{"message" : "dummy_message"}`,
			},
			{ // signal through a valid channel with a unregisted channel
				state: struct {
					channels []string
					level    Level
				}{
					channels: []string{"dummy_channel"},
					level:    WARNING,
				},
				call: struct {
					level   Level
					channel string
					fields  F
					message string
				}{
					level:   FATAL,
					channel: "not_a_valid_dummy_channel",
					fields:  F{},
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
			formatter := NewMockFormatter(ctrl)
			formatter.EXPECT().Format(scn.call.level, scn.call.message, scn.call.fields).Return(scn.expected).Times(scn.callTimes)
			stream, _ := NewFileStream(writer, formatter, scn.state.channels, scn.state.level)
			defer stream.Close()

			if err := stream.Signal(scn.call.channel, scn.call.level, scn.call.message, scn.call.fields); err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}

func Test_FileStream_Broadcast(t *testing.T) {
	t.Run("broadcast message to the writer", func(t *testing.T) {
		scenarios := []struct {
			state struct {
				channels []string
				level    Level
			}
			call struct {
				level   Level
				fields  F
				message string
			}
			callTimes int
			expected  string
		}{
			{ // broadcast through a valid channel with a not filtered level
				state: struct {
					channels []string
					level    Level
				}{
					channels: []string{"dummy_channel"},
					level:    WARNING,
				},
				call: struct {
					level   Level
					fields  F
					message string
				}{
					level:   FATAL,
					fields:  F{},
					message: "dummy_message",
				},
				callTimes: 1,
				expected:  `{"message" : "dummy_message"}`,
			},
			{ // broadcast through a valid channel with a filtered level
				state: struct {
					channels []string
					level    Level
				}{
					channels: []string{"dummy_channel"},
					level:    WARNING,
				},
				call: struct {
					level   Level
					fields  F
					message string
				}{
					level:   DEBUG,
					fields:  F{},
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
			formatter := NewMockFormatter(ctrl)
			formatter.EXPECT().Format(scn.call.level, scn.call.message, scn.call.fields).Return(scn.expected).Times(scn.callTimes)
			stream, _ := NewFileStream(writer, formatter, scn.state.channels, scn.state.level)
			defer stream.Close()

			if err := stream.Broadcast(scn.call.level, scn.call.message, scn.call.fields); err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}
