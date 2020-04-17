package log

import (
	"regexp"
	"testing"
)

func Test_NewJSONFormatter(t *testing.T) {
	t.Run("creates a new json formatter", func(t *testing.T) {
		if NewJSONFormatter() == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_JSONFormatter_Format(t *testing.T) {
	t.Run("correctly format the message", func(t *testing.T) {
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

		formatter := NewJSONFormatter()

		for _, scn := range scenarios {
			result := formatter.Format(scn.level, scn.message, scn.fields)
			matched, _ := regexp.Match(scn.expected, []byte(result))
			if !matched {
				t.Errorf("didn't validate (%s) output", result)
			}
		}
	})
}
