package servlet

// AppProvider is an interface used to define the methods of an object that can
// be registered into a servlet application and register elements in the
// application container and do some necessary boot actions on initialization.
type AppProvider interface {
	Register(*AppContainer) error
	Boot(*AppContainer) error
}
