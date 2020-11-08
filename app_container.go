package servlet

import "fmt"

// AppContainer is a object used to lazy load and store instances of
// registered objects. This is achieved by the registration of factory
// functions that will instantiate the instances as needed.
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
