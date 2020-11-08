package servlet

// LogStreamFactoryStrategy interface defines the methods of the stream
// factory strategy that can validate creation requests and instantiation
// of particular type of stream.
type LogStreamFactoryStrategy interface {
	Accept(sourceType string, args ...interface{}) bool
	AcceptConfig(conf ConfigPartial) bool
	Create(args ...interface{}) (LogStream, error)
	CreateConfig(conf ConfigPartial) (LogStream, error)
}
