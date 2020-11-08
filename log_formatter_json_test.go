package servlet

import (
	"regexp"
	"testing"
)

func Test_NewLogFormatterJSON(t *testing.T) {
	t.Run("new json formatter", func(t *testing.T) {
		if NewLogFormatterJSON() == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_LogFormatterJSON_Format(t *testing.T) {
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

		for _, scn := range scenarios {
			formatter := NewLogFormatterJSON()
			result := formatter.Format(scn.level, scn.message, scn.fields)
			matched, _ := regexp.Match(scn.expected, []byte(result))
			if !matched {
				t.Errorf("didn't validated (%s) output", result)
			}
		}
	})
}
