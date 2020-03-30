package log

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewFileStream(t *testing.T) {
	t.Run("should return nil when missing writer", func(t *testing.T) {
		action := "Creating a new file stream without a writer reference"

		expected := "Invalid nil 'writer' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatter := NewMockFormatter(ctrl)

		stream, err := NewFileStream(nil, formatter, []string{}, WARNING)

		if stream != nil {
			stream.Close()
			t.Errorf("%s returned a valid file stream reference, expected nil", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return nil when missing formatter", func(t *testing.T) {
		action := "Creating a new file stream without a formatter reference"

		expected := "Invalid nil 'formatter' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		writer := NewMockWriter(ctrl)

		stream, err := NewFileStream(writer, nil, []string{}, WARNING)

		if stream != nil {
			stream.Close()
			t.Errorf("%s returned a valid file stream reference, expected nil", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("creates a new file stream", func(t *testing.T) {
		action := "Creating a new file stream"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		writer := NewMockWriter(ctrl)
		writer.EXPECT().Close().Times(1)

		formatter := NewMockFormatter(ctrl)

		stream, err := NewFileStream(writer, formatter, []string{}, WARNING)

		if stream == nil {
			t.Errorf("%s didn't return a valid reference to a new file stream", action)
		} else {
			stream.Close()
		}
		if err != nil {
			t.Errorf("%s returned a unexpected error : %v", action, err)
		}
	})
}

func Test_FileStream_Close(t *testing.T) {
	t.Run("should call the close on the writer only once", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		writer := NewMockWriter(ctrl)
		writer.EXPECT().Close().Times(1)

		formatter := NewMockFormatter(ctrl)

		stream, _ := NewFileStream(writer, formatter, []string{}, WARNING)
		stream.Close()
		stream.Close()
	})
}

func Test_FileStream_Signal(t *testing.T) {
	t.Run("should correctly send the signal message to the writer (if not filtered out)", func(t *testing.T) {
		action := "Calling a signal type of message to the file stream"

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
					channels: []string{"__dummy_channel__"},
					level:    WARNING,
				},
				call: struct {
					level   Level
					channel string
					fields  F
					message string
				}{
					level:   FATAL,
					channel: "__dummy_channel__",
					fields:  F{},
					message: "__dummy_message__",
				},
				callTimes: 1,
				expected:  `{"message" : "__dummy_message__"}`,
			},
			{ // signal through a valid channel with a filtered level
				state: struct {
					channels []string
					level    Level
				}{
					channels: []string{"__dummy_channel__"},
					level:    WARNING,
				},
				call: struct {
					level   Level
					channel string
					fields  F
					message string
				}{
					level:   DEBUG,
					channel: "__dummy_channel__",
					fields:  F{},
					message: "__dummy_message__",
				},
				callTimes: 0,
				expected:  `{"message" : "__dummy_message__"}`,
			},
			{ // signal through a valid channel with a unregisted channel
				state: struct {
					channels []string
					level    Level
				}{
					channels: []string{"__dummy_channel__"},
					level:    WARNING,
				},
				call: struct {
					level   Level
					channel string
					fields  F
					message string
				}{
					level:   FATAL,
					channel: "__not_a_valid_dummy_channel__",
					fields:  F{},
					message: "__dummy_message__",
				},
				callTimes: 0,
				expected:  `{"message" : "__dummy_message__"}`,
			},
		}

		for _, scn := range scenarios {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			writer := NewMockWriter(ctrl)
			writer.EXPECT().Close().Times(1)
			writer.EXPECT().Write([]byte(scn.expected + "\n")).Times(scn.callTimes)

			formatter := NewMockFormatter(ctrl)
			formatter.EXPECT().Format(scn.call.level, scn.call.fields, scn.call.message).Return(scn.expected).Times(scn.callTimes)

			stream, _ := NewFileStream(writer, formatter, scn.state.channels, scn.state.level)
			defer stream.Close()

			if result := stream.Signal(scn.call.channel, scn.call.level, scn.call.fields, scn.call.message); result != nil {
				t.Errorf("%s returned unexpected error (%v)", action, result)
			}
		}
	})
}

func Test_FileStream_Broadcast(t *testing.T) {
	t.Run("should correctly send the broadcast message to the writer (if not filtered out)", func(t *testing.T) {
		action := "Calling a broadcast type of message to the file stream"

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
					channels: []string{"__dummy_channel__"},
					level:    WARNING,
				},
				call: struct {
					level   Level
					fields  F
					message string
				}{
					level:   FATAL,
					fields:  F{},
					message: "__dummy_message__",
				},
				callTimes: 1,
				expected:  `{"message" : "__dummy_message__"}`,
			},
			{ // broadcast through a valid channel with a filtered level
				state: struct {
					channels []string
					level    Level
				}{
					channels: []string{"__dummy_channel__"},
					level:    WARNING,
				},
				call: struct {
					level   Level
					fields  F
					message string
				}{
					level:   DEBUG,
					fields:  F{},
					message: "__dummy_message__",
				},
				callTimes: 0,
				expected:  `{"message" : "__dummy_message__"}`,
			},
		}

		for _, scn := range scenarios {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			writer := NewMockWriter(ctrl)
			writer.EXPECT().Close().Times(1)
			writer.EXPECT().Write([]byte(scn.expected + "\n")).Times(scn.callTimes)

			formatter := NewMockFormatter(ctrl)
			formatter.EXPECT().Format(scn.call.level, scn.call.fields, scn.call.message).Return(scn.expected).Times(scn.callTimes)

			stream, _ := NewFileStream(writer, formatter, scn.state.channels, scn.state.level)
			defer stream.Close()

			if result := stream.Broadcast(scn.call.level, scn.call.fields, scn.call.message); result != nil {
				t.Errorf("%s returned unexpected error (%v)", action, result)
			}
		}
	})
}
