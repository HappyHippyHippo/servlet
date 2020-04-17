package servlet

// Factory is a callback function used to instantiate an object used by
// the application container when a not yet instantiated object is requested.
type Factory func(Container) interface{}
