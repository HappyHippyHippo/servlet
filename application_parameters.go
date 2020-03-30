package servlet

// ContainerEngineID defines the default id used to register the application
// gin engine instance in the application container.
const ContainerEngineID = "servlet.engine"

// ApplicationParameters interface defines the method to interact
// with the application parameters object used to initialize the application.
type ApplicationParameters interface {
	GetEngineID() string
	SetEngineID(ID string) ApplicationParameters
}

type applicationParameters struct {
	engineID string
}

// NewApplicationParameters will instantiate a new application parameters
// object with the requested values.
func NewApplicationParameters(engineID string) ApplicationParameters {
	return &applicationParameters{
		engineID: engineID,
	}
}

// NewDefaultApplicationParameters instantiate a new application
// parameters object with the default values.
func NewDefaultApplicationParameters() ApplicationParameters {
	return &applicationParameters{
		engineID: ContainerEngineID,
	}
}

// GetEngineID getter of the ID parameter value.
func (p applicationParameters) GetEngineID() string {
	return p.engineID
}

// SetEngineID setter of the ID parameter value.
func (p *applicationParameters) SetEngineID(id string) ApplicationParameters {
	p.engineID = id
	return p
}
