package servlet

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Application interface used to define the methods of a servlet application.
type Application interface {
	Boot()
	Engine() Engine
	GetContainer() Container
	SetContainer(container Container) error
	AddProvider(provider Provider) error
	Run(addr ...string) error
}

type application struct {
	engine    Engine
	container Container
	providers []Provider
	boot      bool
}

// NewApplication used to instanciate a new application.
func NewApplication(params Parameters) Application {
	engine := gin.New()
	container := NewContainer()

	container.Add(params.EngineID, func(c Container) interface{} {
		return engine
	})

	return &application{
		engine:    engine,
		container: container,
		providers: []Provider{},
		boot:      false,
	}
}

// Boot initialize the application if not initialized yet.
// The initialization of an application is the calling of the register method
// on all providers, after the registration of all objects in the container,
// the boot method of all providers will be executed.
func (a *application) Boot() {
	if !a.boot {
		for _, p := range a.providers {
			p.Boot(a.container)
		}

		a.boot = true
	}
}

// Engine will retrieve the application underlying gin engine.
func (a application) Engine() Engine {
	return a.engine
}

// Container will retrieve the application underlying container.
func (a application) GetContainer() Container {
	return a.container
}

func (a *application) SetContainer(container Container) error {
	if container == nil {
		return fmt.Errorf("Invalid nil 'container' argument")
	}

	a.container = container
	return nil
}

// AddProvider will register a new provider into the application used
// on the application boot.
func (a *application) AddProvider(provider Provider) error {
	if provider == nil {
		return fmt.Errorf("Invalid nil 'provider' argument")
	}

	a.providers = append(a.providers, provider)
	provider.Register(a.container)

	return nil
}

// Run method will boot the application, if not yet, and the start
// the underlying gin server.
func (a application) Run(addr ...string) error {
	a.Boot()

	return a.engine.Run(addr...)
}
