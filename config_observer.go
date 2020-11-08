package servlet

// ConfigObserver callback function used to be called when a observed
// configuration path has changed.
type ConfigObserver func(interface{}, interface{})
