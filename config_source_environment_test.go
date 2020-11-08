package servlet

import (
	"os"
	"reflect"
	"testing"
)

func Test_NewConfigSourceEnvironment(t *testing.T) {
	t.Run("with empty mapping", func(t *testing.T) {
		if source, err := NewConfigSourceEnvironment(map[string]string{}); source == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			defer source.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if source.mutex == nil {
				t.Error("didn't created the access mutex")
			} else if !reflect.DeepEqual(source.partial, ConfigPartial{}) {
				t.Error("didn't loaded the content correctly")
			}
		}
	})

	t.Run("with root mapping", func(t *testing.T) {
		env := "env"
		value := "value"
		mapping := map[string]string{env: "id"}
		expected := ConfigPartial{"id": value}

		_ = os.Setenv(env, value)
		defer func() { _ = os.Setenv(env, "") }()

		if source, err := NewConfigSourceEnvironment(mapping); source == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			defer source.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if source.mutex == nil {
				t.Error("didn't created the access mutex")
			} else if !reflect.DeepEqual(source.partial, expected) {
				t.Error("didn't loaded the content correctly")
			}
		}
	})

	t.Run("with multi-level mapping", func(t *testing.T) {
		env := "env"
		value := "value"
		mapping := map[string]string{env: "root.node"}
		expected := ConfigPartial{"root": ConfigPartial{"node": value}}

		_ = os.Setenv(env, value)
		defer func() { _ = os.Setenv(env, "") }()

		if source, err := NewConfigSourceEnvironment(mapping); source == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			defer source.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if source.mutex == nil {
				t.Error("didn't created the access mutex")
			} else if !reflect.DeepEqual(source.partial, expected) {
				t.Error("didn't loaded the content correctly")
			}
		}
	})

	t.Run("with multi-level mapping and node override", func(t *testing.T) {
		mapping := map[string]string{
			"env1": "root",
			"env2": "root.node",
		}

		expected := ConfigPartial{"root": ConfigPartial{"node": "value2"}}

		_ = os.Setenv("env1", "value1")
		_ = os.Setenv("env2", "value2")
		defer func() { _ = os.Setenv("env1", ""); _ = os.Setenv("env2", "") }()

		if source, err := NewConfigSourceEnvironment(mapping); source == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			defer source.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if source.mutex == nil {
				t.Error("didn't created the access mutex")
			} else if !reflect.DeepEqual(source.partial, expected) {
				t.Error("didn't loaded the content correctly")
			}
		}
	})
}
