package config

// SourceFactoryStrategy interface defines the methods of the source factory
// strategy that will be used instantiate a particular source type.
type SourceFactoryStrategy interface {
	Accept(stype string, args ...interface{}) bool
	AcceptConfig(conf Partial) bool
	Create(args ...interface{}) (Source, error)
	CreateConfig(conf Partial) (Source, error)
}
