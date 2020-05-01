package servlet

import (
	"fmt"
)

/// ---------------------------------------------------------------------------
/// Closable
/// ---------------------------------------------------------------------------

// Closable is the interface used to signal the container that
// the element must be closed on removal.
type Closable interface {
	Close()
}

/// ---------------------------------------------------------------------------
/// AppContainerFactory
/// ---------------------------------------------------------------------------

// AppContainerFactory is a callback function used to instantiate an object used by
// the application container when a not yet instantiated object is requested.
type AppContainerFactory func(*AppContainer) (interface{}, error)

/// ---------------------------------------------------------------------------
/// AppContainer
/// ---------------------------------------------------------------------------

// AppContainer is a object used to lazy load and store instances of
// registered objects. This is achieved by the registration of factory functions
// that will instantiate the instances as needed.
type AppContainer struct {
	factories map[string]AppContainerFactory
	entries   map[string]interface{}
}

// NewAppContainer instantiates a new container object.
func NewAppContainer() *AppContainer {
	return &AppContainer{
		factories: map[string]AppContainerFactory{},
		entries:   map[string]interface{}{},
	}
}

// Close clean up the container from all the stored objects.
// If the object has been already instantiated and implements the Closable
// interface, then the Close method will be called upon the removing instance.
func (c *AppContainer) Close() {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	for id := range c.factories {
		c.Remove(id)
	}
}

// Has will check if a object is registered with the requested id.
// This does not mean that is instantiated. The instantiation is just executed
// when the instance is requested for the first time.
func (c AppContainer) Has(id string) bool {
	_, ok := c.factories[id]
	return ok
}

// Add will register the requested object defined by his factory method with
// the requested id value.
// If any object was registered previously with the requested id, then the
// object will be removed by calling the Remove method previously the storing
// of the new object factory.
func (c *AppContainer) Add(id string, factory AppContainerFactory) error {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if factory == nil {
		return fmt.Errorf("invalid nil 'factory' argument")
	}

	c.Remove(id)
	c.factories[id] = factory
	return nil
}

// Remove will eliminate the object from the container.
// If the object has been already instantiated and implements the Closable
// interface, then the Close method will be called on the removing instance.
func (c *AppContainer) Remove(id string) {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if entry, ok := c.entries[id]; ok {
		switch e := entry.(type) {
		case Closable:
			e.Close()
		}
	}
	delete(c.factories, id)
	delete(c.entries, id)
}

// Get will retrieve the requested object from the container.
// If the object has not yet been instantiated, then the factory method will be
// executed to instantiate it.
func (c *AppContainer) Get(id string) (interface{}, error) {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if entry, ok := c.entries[id]; ok {
		return entry, nil
	}

	if factory, ok := c.factories[id]; ok {
		entry, err := factory(c)
		if err != nil {
			return nil, err
		}
		c.entries[id] = entry
		return c.entries[id], nil
	}

	return nil, fmt.Errorf("entry '%s' not registered in the container", id)
}

/// ---------------------------------------------------------------------------
/// AppProvider
/// ---------------------------------------------------------------------------

// AppProvider is an interface used to define the methods of an object that can
// be registered into a servlet application and register elements in the
// application container and do some necessary boot actions on initialization.
type AppProvider interface {
	Register(*AppContainer) error
	Boot(*AppContainer) error
}

/// ---------------------------------------------------------------------------
/// App
/// ---------------------------------------------------------------------------

// App interface used to define the methods of a servlet application.
type App struct {
	container *AppContainer
	providers []AppProvider
	boot      bool
}

// NewApp used to instantiate a new application.
func NewApp() *App {
	return &App{
		container: NewAppContainer(),
		providers: []AppProvider{},
		boot:      false,
	}
}

// Container will retrieve the application underlying container.
func (a App) Container() *AppContainer {
	return a.container
}

// Add will register a new provider into the application used
// on the application boot.
func (a *App) Add(provider AppProvider) error {
	if a == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if provider == nil {
		return fmt.Errorf("invalid nil 'provider' argument")
	}

	a.providers = append(a.providers, provider)
	if err := provider.Register(a.container); err != nil {
		a.providers = a.providers[:len(a.providers)-1]
		return err
	}

	return nil
}

// Boot initialize the application if not initialized yet.
// The initialization of an application is the calling of the register method
// on all providers, after the registration of all objects in the container,
// the boot method of all providers will be executed.
func (a *App) Boot() error {
	if a == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if !a.boot {
		for _, p := range a.providers {
			if err := p.Boot(a.container); err != nil {
				return err
			}
		}

		a.boot = true
	}
	return nil
}
