package log

import (
	"reflect"
	"regexp"
	"testing"
)

type logMessageMatcher struct {
	regex string
}

func (m logMessageMatcher) Matches(x interface{}) bool {
	c, _ := regexp.Compile("\"message\":\"" + m.regex + "\"")
	return c.MatchString(x.(string))
}

func (m logMessageMatcher) String() string {
	return "validate the regex " + m.regex
}

func Test_Stream_HasChannel(t *testing.T) {
	channels := []string{"channel.1", "channel.2"}
	stream := &stream{nil, channels, WARNING}

	t.Run("check the channel registration", func(t *testing.T) {

		if !stream.HasChannel("channel.1") {
			t.Errorf("'channel.1' channel was not found")
		} else if !stream.HasChannel("channel.2") {
			t.Errorf("'channel.2' channel was not found")
		} else if stream.HasChannel("channel.3") {
			t.Errorf("'channel.3' channel was found")
		}
	})
}

func Test_Stream_ListChannels(t *testing.T) {
	channels := []string{"channel.1", "channel.2"}
	stream := &stream{nil, channels, WARNING}

	t.Run("list the registed channels", func(t *testing.T) {
		if result := stream.ListChannels(); !reflect.DeepEqual(result, channels) {
			t.Errorf("returned the (%v) list of channels", result)
		}
	})
}

func Test_Stream_AddChannel(t *testing.T) {
	t.Run("register a new channel", func(t *testing.T) {
		scenarios := []struct {
			state struct {
				channels []string
				level    Level
			}
			channel  string
			expected []string
		}{
			{ // adding into a empty list
				state: struct {
					channels []string
					level    Level
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
					level    Level
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
					level    Level
				}{
					channels: []string{"channel.1", "channel.2", "channel.3"},
					level:    DEBUG,
				},
				channel:  "channel.2",
				expected: []string{"channel.1", "channel.2", "channel.3"},
			},
		}

		for _, scn := range scenarios {
			stream := &stream{nil, scn.state.channels, scn.state.level}
			stream.AddChannel(scn.channel)

			if result := stream.ListChannels(); !reflect.DeepEqual(result, scn.expected) {
				t.Errorf("returned the (%v) list of channels", result)
			}
		}
	})
}

func Test_Stream_RemoveChannel(t *testing.T) {
	t.Run("unregister a channel", func(t *testing.T) {
		scenarios := []struct {
			state struct {
				channels []string
				level    Level
			}
			channel  string
			expected []string
		}{
			{ // removing from an empty list
				state: struct {
					channels []string
					level    Level
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
					level    Level
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
					level    Level
				}{
					channels: []string{"channel.1", "channel.2", "channel.3"},
					level:    DEBUG,
				},
				channel:  "channel.2",
				expected: []string{"channel.1", "channel.3"},
			},
		}

		for _, scn := range scenarios {
			stream := &stream{nil, scn.state.channels, scn.state.level}
			stream.RemoveChannel(scn.channel)

			if result := stream.ListChannels(); !reflect.DeepEqual(result, scn.expected) {
				t.Errorf("returned the (%v) list of channels", result)
			}
		}
	})
}

func Test_Stream_Level(t *testing.T) {
	level := WARNING
	stream := &stream{nil, []string{}, level}

	t.Run("retrieve the filtering level", func(t *testing.T) {
		if result := stream.Level(); result != level {
			t.Errorf("returned the (%v) level", result)
		}
	})
}
