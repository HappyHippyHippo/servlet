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
	t.Run("should correctly check the registration of a channel", func(t *testing.T) {
		action := "Checking for the existence of a channel in the file stream"

		channels := []string{"channel.1", "channel.2"}

		stream := &stream{nil, channels, WARNING}

		if !stream.HasChannel("channel.1") {
			t.Errorf("%s the 'channel.1' channel was not found", action)
		}
		if !stream.HasChannel("channel.2") {
			t.Errorf("%s the 'channel.2' channel was not found", action)
		}
		if stream.HasChannel("channel.3") {
			t.Errorf("%s The 'channel.3' channel was found", action)
		}
	})
}

func Test_Stream_ListChannels(t *testing.T) {
	t.Run("should correctly list the registed channels", func(t *testing.T) {
		action := "Retrieveing the list of channels in the file stream"

		channels := []string{"channel.1", "channel.2"}

		stream := &stream{nil, channels, WARNING}

		if result := stream.ListChannels(); !reflect.DeepEqual(result, channels) {
			t.Errorf("%s retrieved the (%v) list of channels, expected (%v)", action, result, channels)
		}
	})
}

func Test_Stream_AddChannel(t *testing.T) {
	t.Run("should correctly register a new channel", func(t *testing.T) {
		action := "Adding a channel into the file stream"

		scenarios := []struct {
			state struct {
				channels []string
				level    Level
			}
			channel  string
			expected []string
		}{
			// adding into a empty list
			{
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
			// adding should keep sorting
			{
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
			// adding an already existent should result in a no-op
			{
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
				t.Errorf("%s unexpected (%v) resulting list of channels, expected (%v)", action, result, scn.expected)
			}
		}
	})
}

func Test_Stream_RemoveChannel(t *testing.T) {
	t.Run("should correctly unregister a channel", func(t *testing.T) {
		action := "Removing a channel from the file stream"

		scenarios := []struct {
			state struct {
				channels []string
				level    Level
			}
			channel  string
			expected []string
		}{
			// removing from an empty list
			{
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
			// removing a non existing channel
			{
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
			// removing an existing channel
			{
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
				t.Errorf("%s unexpected (%v) resulting list of channels, expected (%v)", action, result, scn.expected)
			}
		}
	})
}

func Test_Stream_Level(t *testing.T) {
	t.Run("should correctly retrieve the filtering level", func(t *testing.T) {
		action := "Retrieving the filtering level of the file stream"

		level := WARNING

		stream := &stream{nil, []string{}, level}

		if result := stream.Level(); result != level {
			t.Errorf("%s received the (%v) list of channels, expected (%v)", action, result, level)
		}
	})
}
