package servlet

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_NewConfigDecoderFactoryStrategyYaml(t *testing.T) {
	t.Run("new strategy", func(t *testing.T) {
		if strategy := NewConfigDecoderFactoryStrategyYaml(); strategy == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_ConfigDecoderFactoryStrategyYaml_Accept(t *testing.T) {
	t.Run("accept only yaml format", func(t *testing.T) {
		scenarios := []struct {
			format   string
			expected bool
		}{
			{ // test yaml format
				format:   ConfigDecoderFormatYAML,
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

			reader := NewMockReader(ctrl)
			strategy := NewConfigDecoderFactoryStrategyYaml()

			if check := strategy.Accept(scn.format, reader); check != scn.expected {
				t.Errorf("returned (%v) when checking (%s) format", check, scn.format)
			}
		}
	})

	t.Run("no extra arguments", func(t *testing.T) {
		strategy := NewConfigDecoderFactoryStrategyYaml()
		if strategy.Accept(ConfigDecoderFormatYAML) {
			t.Error("returned true")
		}
	})

	t.Run("first extra argument is not a io.Reader interface", func(t *testing.T) {
		strategy := NewConfigDecoderFactoryStrategyYaml()
		if strategy.Accept(ConfigDecoderFormatYAML, "string") {
			t.Error("returned true")
		}
	})
}

func Test_ConfigDecoderFactoryStrategyYaml_Create(t *testing.T) {
	t.Run("create the decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockReader(ctrl)
		strategy := NewConfigDecoderFactoryStrategyYaml()

		if decoder, err := strategy.Create(reader); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if decoder == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch decoder.(type) {
			case *ConfigDecoderYaml:
			default:
				t.Error("didn't returned a YAML decoder")
			}
		}
	})
}
