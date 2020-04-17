package servlet

import "os"

const (
	// ContainerEngineID defines the default id used to register the
	// application gin engine instance in the application container.
	ContainerEngineID = "servlet.engine"

	// EnvContainerEngineID defines the environment variable used to
	// override the default value for the container gin engine id
	EnvContainerEngineID = "SERVLET_CONTAINER_ENGINE_ID"
)

// Parameters defines the application parameters storing structure
// that will be needed when instantiating a new application
type Parameters struct {
	EngineID string
}

// NewParameters will instantiate a new application parameters
// object with the requested values.
func NewParameters(engineID string) Parameters {
	return Parameters{
		EngineID: engineID,
	}
}

// NewDefaultParameters instantiate a new application
// parameters object with the default values.
func NewDefaultParameters() Parameters {
	engineID := ContainerEngineID
	if env := os.Getenv(EnvContainerEngineID); env != "" {
		engineID = env
	}

	return Parameters{
		EngineID: engineID,
	}
}
