package log

// FormatterFactoryStrategy interface defines the methods of the formatter
// factory strategy that can instantiate a particular formatter.
type FormatterFactoryStrategy interface {
	Accept(format string, args ...interface{}) bool
	Create(args ...interface{}) (Formatter, error)
}
