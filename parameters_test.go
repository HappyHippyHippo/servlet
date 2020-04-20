package servlet

import (
	"os"
	"testing"
)

func Test_NewParameters(t *testing.T) {
	id := "id"
	parameters := NewParameters(id)

	t.Run("creates a new parameters", func(t *testing.T) {
		if value := parameters.EngineID; value != id {
			t.Errorf("stored (%v) engine ID", value)
		}
	})
}

func Test_NewDefaultParameters(t *testing.T) {
	t.Run("creates a new parameters with the default values", func(t *testing.T) {
		parameters := NewDefaultParameters()
		if value := parameters.EngineID; value != ContainerEngineID {
			t.Errorf("stored (%v) engine ID", value)
		}
	})

	t.Run("creates a new parameters with the engine id environment override", func(t *testing.T) {
		engineID := "id"
		os.Setenv(EnvContainerEngineID, engineID)
		defer os.Setenv(EnvContainerEngineID, "")

		parameters := NewDefaultParameters()
		if value := parameters.EngineID; value != engineID {
			t.Errorf("stored (%v) engine ID", value)
		}
	})
}
