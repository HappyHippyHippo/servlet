package servlet

// ConfigSourceObservable interface extends the Source interface with methods
// specific to sources that will be checked for updates in a regular
// periodicity defined in the config object where the source will be
// registered.
type ConfigSourceObservable interface {
	ConfigSource
	Reload() (bool, error)
}
