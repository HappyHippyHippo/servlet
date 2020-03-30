package servlet

// Provider is an interface used to define the method of an object that can
// be registed into a servlet application and register elements in such
// application.
type Provider interface {
	Register(container Container)
	Boot(container Container)
}
