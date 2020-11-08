package servlet

// LogFormatterFactoryStrategyJSON defines the log formatter instantiation
// strategy to be registered in the factory so a Json based log formatter
// could be instantiated.
type LogFormatterFactoryStrategyJSON struct{}

// NewLogFormatterFactoryStrategyJSON instantiate a new json logging output
// formatter factory strategy that will enable the formatter factory to
// instantiate a new content to json formatter.
func NewLogFormatterFactoryStrategyJSON() *LogFormatterFactoryStrategyJSON {
	return &LogFormatterFactoryStrategyJSON{}
}

// Accept will check if the formatter factory strategy can instantiate a
// formatter of the requested format.
func (LogFormatterFactoryStrategyJSON) Accept(format string, _ ...interface{}) bool {
	return format == LogFormatterFormatJSON
}

// Create will instantiate the desired formatter instance.
func (LogFormatterFactoryStrategyJSON) Create(_ ...interface{}) (LogFormatter, error) {
	return NewLogFormatterJSON(), nil
}
