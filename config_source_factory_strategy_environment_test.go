package servlet

import (
	"github.com/golang/mock/gomock"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_NewConfigSourceFactoryStrategyEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("new env source factory strategy", func(t *testing.T) {
		if strategy, err := NewConfigSourceFactoryStrategyEnvironment(); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if strategy == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_ConfigSourceFactoryStrategyEnvironment_Accept(t *testing.T) {
	t.Run("don't accept if at least 1 extra arguments are passed", func(t *testing.T) {
		sourceType := ConfigSourceTypeEnv

		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()
		if strategy.Accept(sourceType) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if the mappings is not a string", func(t *testing.T) {
		sourceType := ConfigSourceTypeEnv

		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()
		if strategy.Accept(sourceType, 1) {
			t.Error("returned true")
		}
	})

	t.Run("accept only env type", func(t *testing.T) {
		scenarios := []struct {
			sourceType string
			expected   bool
		}{
			{ // test env type
				sourceType: ConfigSourceTypeEnv,
				expected:   true,
			},
			{ // test non-file type (file)
				sourceType: ConfigSourceTypeFile,
				expected:   false,
			},
		}

		for _, scn := range scenarios {
			mapping := map[string]string{}

			strategy, _ := NewConfigSourceFactoryStrategyEnvironment()
			if check := strategy.Accept(scn.sourceType, mapping); check != scn.expected {
				t.Errorf("for the type (%s), returned (%v)", scn.sourceType, check)
			}
		}
	})
}

func Test_ConfigSourceFactoryStrategyEnvironment_AcceptConfig(t *testing.T) {
	t.Run("don't accept if type is missing", func(t *testing.T) {
		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()

		partial := ConfigPartial{}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if type is not a string", func(t *testing.T) {
		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()

		partial := ConfigPartial{"type": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if mapping is missing", func(t *testing.T) {
		sourceType := ConfigSourceTypeEnv

		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()

		partial := ConfigPartial{"type": sourceType}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if mapping is not a string", func(t *testing.T) {
		sourceType := ConfigSourceTypeEnv

		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()

		partial := ConfigPartial{"type": sourceType, "mapping": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if invalid type", func(t *testing.T) {
		mapping := map[string]string{}

		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()

		partial := ConfigPartial{"type": ConfigSourceTypeFile, "mapping": mapping}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("accept config", func(t *testing.T) {
		sourceType := ConfigSourceTypeEnv
		mapping := map[string]string{}

		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()

		partial := ConfigPartial{"type": sourceType, "mapping": mapping}
		if !strategy.AcceptConfig(partial) {
			t.Error("returned false")
		}
	})
}

func Test_CConfigSourceFactoryStrategyEnvironment_Create(t *testing.T) {
	t.Run("non-map mapping", func(t *testing.T) {
		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()

		if source, err := strategy.Create(123); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the source", func(t *testing.T) {
		env := "env"
		path := "root"
		value := "value"
		mapping := map[string]string{env: path}
		expected := ConfigPartial{path: value}

		_ = os.Setenv(env, value)
		defer func() { _ = os.Setenv(env, "") }()

		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()

		if source, err := strategy.Create(mapping); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if source == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch s := source.(type) {
			case *ConfigSourceEnvironment:
				if !reflect.DeepEqual(s.partial, expected) {
					t.Error("didn't loaded the content correctly")
				}
			default:
				t.Error("didn't returned a new env source")
			}
		}
	})
}

func Test_ConfigSourceFactoryStrategyEnvironment_CreateConfig(t *testing.T) {
	t.Run("non-map mapping", func(t *testing.T) {
		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()

		conf := ConfigPartial{"mapping": 123}
		if source, err := strategy.CreateConfig(conf); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the source", func(t *testing.T) {
		env := "env"
		path := "root"
		value := "value"
		expected := ConfigPartial{path: value}

		_ = os.Setenv(env, value)
		defer func() { _ = os.Setenv(env, "") }()

		strategy, _ := NewConfigSourceFactoryStrategyEnvironment()

		conf := ConfigPartial{"mapping": ConfigPartial{env: path}}

		if source, err := strategy.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if source == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch s := source.(type) {
			case *ConfigSourceEnvironment:
				if !reflect.DeepEqual(s.partial, expected) {
					t.Error("didn't loaded the content correctly")
				}
			default:
				t.Error("didn't returned a new env source")
			}
		}
	})
}
