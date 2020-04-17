package servlet

// Provider is an interface used to define the methods of an object that can
// be registed into a servlet application and register elements in the
// application container and do some necessary boot actions on initialization.
type Provider interface {
	Register(container Container)
	Boot(container Container)
}
