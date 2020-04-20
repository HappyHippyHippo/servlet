package config

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewYamlDecoderFactoryStrategy(t *testing.T) {
	t.Run("creates a new strategy", func(t *testing.T) {
		if strategy := NewYamlDecoderFactoryStrategy(); strategy == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_YamlDecoderFactoryStrategy_Accept(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockReader(ctrl, "{}")
	strategy := NewYamlDecoderFactoryStrategy()

	t.Run("accept only yaml format", func(t *testing.T) {
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
			if check := strategy.Accept(scn.format, reader); check != scn.expected {
				t.Errorf("returned (%v) when checking (%s) format", check, scn.format)
			}
		}
	})

	t.Run("no extra arguments", func(t *testing.T) {
		if strategy.Accept(DecoderFormatYAML) {
			t.Errorf("returned true")
		}
	})

	t.Run("first extra argument is not a io.Reader interface", func(t *testing.T) {
		if strategy.Accept(DecoderFormatYAML, "string") {
			t.Errorf("returned true")
		}
	})
}

func Test_YamlDecoderFactoryStrategy_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockReader(ctrl, "{}")
	strategy := NewYamlDecoderFactoryStrategy()

	t.Run("create the decoder", func(t *testing.T) {
		if decoder, err := strategy.Create(reader); err != nil {
			t.Errorf("return the (%v) error", err)
		} else if decoder == nil {
			t.Errorf("didn't return a valid reference")
		} else {
			switch decoder.(type) {
			case *yamlDecoder:
			default:
				t.Errorf("didn't return a YAML decoder")
			}
		}
	})
}
