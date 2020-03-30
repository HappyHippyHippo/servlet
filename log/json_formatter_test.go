package log

import (
	"regexp"
	"testing"
)

func Test_NewJSONFormatter(t *testing.T) {
	t.Run("creates a new json formatter", func(t *testing.T) {
		action := "Creating a new json formatter"

		formatter := NewJSONFormatter()
		if formatter == nil {
			t.Errorf("%s didn't return a valid reference to a new json formatter", action)
		}
	})
}

func Test_JSONFormatter_Format(t *testing.T) {
	t.Run("should correctly format the message", func(t *testing.T) {
		action := "Formatting the requested message"

		scenarios := []struct {
			level    Level
			fields   F
			message  string
			expected string
		}{
			{ // test level FATAL
				level:    FATAL,
				fields:   nil,
				message:  "",
				expected: `"level"\s*\:\s*"fatal"`,
			},
			{ // test level ERROR
				level:    ERROR,
				fields:   nil,
				message:  "",
				expected: `"level"\s*\:\s*"error"`,
			},
			{ // test level WARNING
				level:    WARNING,
				fields:   nil,
				message:  "",
				expected: `"level"\s*\:\s*"warning"`,
			},
			{ // test level NOTICE
				level:    NOTICE,
				fields:   nil,
				message:  "",
				expected: `"level"\s*\:\s*"notice"`,
			},
			{ // test level DEBUG
				level:    DEBUG,
				fields:   nil,
				message:  "",
				expected: `"level"\s*\:\s*"debug"`,
			},
			{ // test fields (single value)
				level:    DEBUG,
				fields:   F{"field1": "value1"},
				message:  "",
				expected: `"field1"\s*\:\s*"value1"`,
			},
			{ // test fields (multiple value)
				level:    DEBUG,
				fields:   F{"field1": "value1", "field2": "value2"},
				message:  "",
				expected: `"field1"\s*\:\s*"value1"|"field2"\s*\:\s*"value2"`,
			},
			{ // test message
				level:    DEBUG,
				fields:   nil,
				message:  "My message",
				expected: `"message"\s*\:\s*"My\smessage"`,
			},
		}

		for _, scn := range scenarios {
			formatter := NewJSONFormatter()
			result := formatter.Format(scn.level, scn.fields, scn.message)
			matched, _ := regexp.Match(scn.expected, []byte(result))
			if !matched {
				t.Errorf("%s didn't correctly validated agains the (%s) regexp : (%s)", action, scn.expected, result)
			}
		}
	})
}
