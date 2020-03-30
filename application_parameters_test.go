package servlet

import (
	reflect "reflect"
	"testing"
)

func Test_NewApplicationParameters(t *testing.T) {
	t.Run("should creates a new application parameters", func(t *testing.T) {
		action := "Creating the new application parameters"

		id := "__dummy_ID__"

		parameters := NewApplicationParameters(id)

		if check := parameters.GetEngineID(); check != id {
			t.Errorf("%s didn't store the expected (%v) engine ID, returned (%v)", action, id, check)
		}
	})
}

func Test_NewDefaultApplicationParameters(t *testing.T) {
	t.Run("should creates a new application parameters with the default values", func(t *testing.T) {
		action := "Creating the new application parameters without values"

		parameters := NewDefaultApplicationParameters()

		if check := parameters.GetEngineID(); check != ContainerEngineID {
			t.Errorf("%s didn't store the expected (%v) engine ID, returned (%v)", action, ContainerEngineID, check)
		}
	})
}

func Test_ApplicationParameters_GetEngineID(t *testing.T) {
	t.Run("should retrieve the assigned engine ID given at creation", func(t *testing.T) {
		action := "Retrieving the engine ID"

		id := "__dummy_ID__"

		parameters := NewApplicationParameters(id)

		if check := parameters.GetEngineID(); check != id {
			t.Errorf("%s didn't store the expected (%v) engine ID, returned (%v)", action, id, check)
		}
	})
}

func Test_ApplicationParameters_SetEngineID(t *testing.T) {
	t.Run("should assign a new engine ID", func(t *testing.T) {
		action := "Assigning the engine ID"

		id := "__dummy_ID__"

		parameters := NewDefaultApplicationParameters()
		if check := parameters.SetEngineID(id); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetEngineID(); check != id {
			t.Errorf("%s didn't store the expected (%v) engine ID, returned (%v)", action, id, check)
		}
	})
}
