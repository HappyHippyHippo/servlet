package log

import (
	"testing"
)

func Test_NewJSONFormatterFactoryStrategy(t *testing.T) {
	t.Run("creates a new json formatter factory strategy", func(t *testing.T) {
		if NewJSONFormatterFactoryStrategy() == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_JSONFormatterFactoryStrategy_Accept(t *testing.T) {
	t.Run("accept only json format", func(t *testing.T) {
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

		strategy := NewJSONFormatterFactoryStrategy()

		for _, scn := range scenarios {
			if check := strategy.Accept(scn.format); check != scn.expected {
				t.Errorf("returned (%v) for the (%s) format", check, scn.format)
			}
		}
	})
}

func Test_JSONFormatterFactoryStrategy_Create(t *testing.T) {
	strategy := NewJSONFormatterFactoryStrategy()

	t.Run("create the requested json formatter", func(t *testing.T) {
		if formatter, err := strategy.Create(); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if formatter == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			switch formatter.(type) {
			case *jsonFormatter:
			default:
				t.Errorf("didn't return a new json formatter")
			}
		}
	})
}
