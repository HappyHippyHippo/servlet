package servlet

import "fmt"

// Factory is a callback function used to instantiate an object
// in the container.
type Factory func(Container) interface{}

// Container object interface defines the methods of a container object used
// to lazy load and store instances of registered objects.
type Container interface {
	Close() error
	Has(id string) bool
	Add(id string, factory Factory) error
	Remove(id string)
	Get(id string) interface{}
}

// Container is a object used to lazy load and store instances of
// registered objects.
type container struct {
	factories map[string]Factory
	entries   map[string]interface{}
}

// NewContainer instanciates a new container object.
func NewContainer() Container {
	return &container{
		factories: make(map[string]Factory),
		entries:   make(map[string]interface{}),
	}
}

// Close clean up the container from all the stored objects.
// If the object has been already instanciated and implements the Closable
// interface, then the Close method will be called apon the removing instance.
func (c *container) Close() error {
	for id := range c.factories {
		c.Remove(id)
	}
	return nil
}

// Has will check if a object is registed with the requested id.
func (c container) Has(id string) bool {
	_, ok := c.factories[id]
	return ok
}

// Add will register the requested object defined nby his factory method with
// the requested id value.
// If any object was registed previously with the requested id, then the
// object will be removed by calling the Remove method previously the storing
// of the new object factory.
func (c *container) Add(id string, factory Factory) error {
	if factory == nil {
		return fmt.Errorf("Invalid nil 'factory' argument")
	}

	c.Remove(id)
	c.factories[id] = factory
	return nil
}

// Remove will eliminate the object from the container.
// If the object has been already instanciated and implements the Closable
// interface, then the Close method will be called apon the removing instance.
func (c *container) Remove(id string) {
	if entry, ok := c.entries[id]; ok {
		switch entry.(type) {
		case Closable:
			entry.(Closable).Close()
		}
	}
	delete(c.factories, id)
	delete(c.entries, id)
}

// Get will retrieve the requested object from the container.
// If the object has not yet been instanciated, then the factory method will be
// executed to instanciate it.
func (c *container) Get(id string) interface{} {
	if entry, ok := c.entries[id]; ok {
		return entry
	}

	if factory, ok := c.factories[id]; ok {
		c.entries[id] = factory(c)
		return c.entries[id]
	}

	panic(fmt.Errorf("Object '%s' not registed in the container", id))
}
