package servlet

import "testing"

func Test_NewLogFormatterFactoryStrategyJSON(t *testing.T) {
	t.Run("new json formatter factory strategy", func(t *testing.T) {
		if NewLogFormatterFactoryStrategyJSON() == nil {
			t.Errorf("didn't returned a valid reference")
		}
	})
}

func Test_LogFormatterFactoryStrategyJSON_Accept(t *testing.T) {
	t.Run("accept only json format", func(t *testing.T) {
		scenarios := []struct {
			format   string
			expected bool
		}{
			{ // test json format
				format:   LogFormatterFormatJSON,
				expected: true,
			},
			{ // test non-json format (yaml)
				format:   "yaml",
				expected: false,
			},
		}

		strategy := NewLogFormatterFactoryStrategyJSON()

		for _, scn := range scenarios {
			if check := strategy.Accept(scn.format); check != scn.expected {
				t.Errorf("returned (%v) for the (%s) format", check, scn.format)
			}
		}
	})
}

func Test_LogFormatterFactoryStrategyJSON_Create(t *testing.T) {
	t.Run("create json formatter", func(t *testing.T) {
		strategy := NewLogFormatterFactoryStrategyJSON()
		if formatter, err := strategy.Create(); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if formatter == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch formatter.(type) {
			case *LogFormatterJSON:
			default:
				t.Errorf("didn't returned a new json formatter")
			}
		}
	})
}
