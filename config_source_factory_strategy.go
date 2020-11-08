package servlet

// ConfigSourceFactoryStrategy interface defines the methods of the source
// factory strategy that will be used instantiate a particular source type.
type ConfigSourceFactoryStrategy interface {
	Accept(sourceType string, args ...interface{}) bool
	AcceptConfig(conf ConfigPartial) bool
	Create(args ...interface{}) (ConfigSource, error)
	CreateConfig(conf ConfigPartial) (ConfigSource, error)
}
