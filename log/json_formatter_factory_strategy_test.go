package log

import (
	"testing"
)

func Test_NewJSONFormatterFactoryStrategy(t *testing.T) {
	t.Run("creates a new json formatter factory strategy", func(t *testing.T) {
		action := "Creating a json formatter factory strategy"

		strategy := NewJSONFormatterFactoryStrategy()

		if strategy == nil {
			t.Errorf("%s didn't return a valid reference to a new json formatter factory strategy", action)
		}
	})
}

func Test_JSONFormatterFactoryStrategy_Accept(t *testing.T) {
	t.Run("should accept only json format", func(t *testing.T) {
		action := "Checking the accepting format"

		scenarios := []struct {
			format   string
			expected bool
		}{
			{ // test json format
				format:   FormatterFormatJSON,
				expected: true,
			},
			{ // test non-json format (yaml)
				format:   "yaml",
				expected: false,
			},
		}

		for _, scn := range scenarios {
			strategy := NewJSONFormatterFactoryStrategy()

			if check := strategy.Accept(scn.format); check != scn.expected {
				t.Errorf("%s didn't returned the expected (%v) for the format (%s), returned (%v)", action, scn.expected, scn.format, check)
			}
		}
	})
}

func Test_JSONFormatterFactoryStrategy_Create(t *testing.T) {
	t.Run("should create the requested json formatter", func(t *testing.T) {
		action := "Creating a new json formatter"

		strategy := NewJSONFormatterFactoryStrategy()

		formatter, err := strategy.Create()
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if formatter == nil {
			t.Errorf("%s didn't returned a valid log formatter reference", action)
		} else {
			switch formatter.(type) {
			case *jsonFormatter:
			default:
				t.Errorf("%s didn't return a valid reference to a new json formatter reference", action)
			}
		}
	})
}
