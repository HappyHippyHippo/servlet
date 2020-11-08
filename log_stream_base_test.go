package servlet

import (
	"reflect"
	"testing"
)

func Test_LogStreamBase_HasChannel(t *testing.T) {
	channels := []string{"channel.1", "channel.2"}
	stream := &LogStreamBase{nil, channels, WARNING}

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

func Test_LogStreamBase_ListChannels(t *testing.T) {
	channels := []string{"channel.1", "channel.2"}
	stream := &LogStreamBase{nil, channels, WARNING}

	t.Run("list the registered channels", func(t *testing.T) {
		if result := stream.ListChannels(); !reflect.DeepEqual(result, channels) {
			t.Errorf("returned the (%v) list of channels", result)
		}
	})
}

func Test_LogStreamBase_AddChannel(t *testing.T) {
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

		var stream *LogStreamBase
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
			stream := &LogStreamBase{nil, scn.state.channels, scn.state.level}
			stream.AddChannel(scn.channel)

			if result := stream.ListChannels(); !reflect.DeepEqual(result, scn.expected) {
				t.Errorf("returned the (%v) list of channels", result)
			}
		}
	})
}

func Test_LogStreamBase_RemoveChannel(t *testing.T) {
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

		var stream *LogStreamBase
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
			stream := &LogStreamBase{nil, scn.state.channels, scn.state.level}
			stream.RemoveChannel(scn.channel)

			if result := stream.ListChannels(); !reflect.DeepEqual(result, scn.expected) {
				t.Errorf("returned the (%v) list of channels", result)
			}
		}
	})
}

func Test_LogStreamBase_Level(t *testing.T) {
	level := WARNING
	stream := &LogStreamBase{nil, []string{}, level}

	t.Run("retrieve the filtering level", func(t *testing.T) {
		if result := stream.Level(); result != level {
			t.Errorf("returned the (%v) level", result)
		}
	})
}
