package servlet

// LogFormatterFactoryStrategy interface defines the methods of the formatter
// factory strategy that can validate creation requests and instantiation
// of particular decoder.
type LogFormatterFactoryStrategy interface {
	Accept(format string, args ...interface{}) bool
	Create(args ...interface{}) (LogFormatter, error)
}
