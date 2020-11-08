package servlet

import (
	"fmt"
)

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
