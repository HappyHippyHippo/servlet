package log

// FormatterFactoryStrategy interface defines the methods of the formatter
// factory strategy that can validate creation requests and instantiation
// of particular decoder.
type FormatterFactoryStrategy interface {
	Accept(format string, args ...interface{}) bool
	Create(args ...interface{}) (Formatter, error)
}
