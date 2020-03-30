package config

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewYamlDecoderFactoryStrategy(t *testing.T) {
	t.Run("creates a new yaml decoder factory strategy", func(t *testing.T) {
		action := "Creating a yaml decoder factory strategy"

		strategy := NewYamlDecoderFactoryStrategy()

		if strategy == nil {
			t.Errorf("%s didn't return a valid reference to a new yaml decoder factory strategy", action)
		}
	})
}

func Test_YamlDecoderFactoryStrategy_Accept(t *testing.T) {
	t.Run("should accept only yaml format", func(t *testing.T) {
		action := "Checking the accepting format"

		scenarios := []struct {
			format   string
			expected bool
		}{
			{ // test yaml format
				format:   DecoderFormatYAML,
				expected: true,
			},
			{ // test non-yaml format (json)
				format:   "json",
				expected: false,
			},
		}

		for _, scn := range scenarios {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			reader := NewMockReader(ctrl, "{}")

			strategy := NewYamlDecoderFactoryStrategy()

			if check := strategy.Accept(scn.format, reader); check != scn.expected {
				t.Errorf("%s didn't returned the expected (%v) for the format (%s), returned (%v)", action, scn.expected, scn.format, check)
			}
		}
	})

	t.Run("should not accept if no extra arguments are passed (the reader)", func(t *testing.T) {
		action := "Checking the acceptance with no extra arguments"

		strategy := NewYamlDecoderFactoryStrategy()

		if strategy.Accept(DecoderFormatYAML) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if the first extra argument is not a io.Reader interface", func(t *testing.T) {
		action := "Checking the acceptance when the first extra argument not a io.Reader interface"

		strategy := NewYamlDecoderFactoryStrategy()

		if strategy.Accept(DecoderFormatYAML, "__string__") {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})
}

func Test_YamlDecoderFactoryStrategy_Create(t *testing.T) {
	t.Run("should create the requested yaml config decoder", func(t *testing.T) {
		action := "Creating a new yaml decoder"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockReader(ctrl, "{}")

		strategy := NewYamlDecoderFactoryStrategy()

		decoder, err := strategy.Create(reader)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if decoder == nil {
			t.Errorf("%s didn't returned a valid config decoder reference", action)
		} else {
			switch decoder.(type) {
			case *yamlDecoder:
			default:
				t.Errorf("%s didn't return a valid reference to a new YAML decoder reference", action)
			}
		}
	})
}
