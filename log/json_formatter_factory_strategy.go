package log

const (
	// FormatterFormatJSON defines the value to be used to declare a JSON
	// log formatter format.
	FormatterFormatJSON = "json"
)

type jsonFormatterFactoryStrategy struct{}

// NewJSONFormatterFactoryStrategy instantiate a new yaml decoder factory
// strategy that will enable the decoder factory to instantiate a new yaml
// decoder.
func NewJSONFormatterFactoryStrategy() FormatterFactoryStrategy {
	return &jsonFormatterFactoryStrategy{}
}

// Accept will check if the formatter factory strategy can instantiate a
// formatter of the requested format.
func (jsonFormatterFactoryStrategy) Accept(format string, args ...interface{}) bool {
	return format == FormatterFormatJSON
}

// Create will instantiate the desired formatter instance.
func (jsonFormatterFactoryStrategy) Create(args ...interface{}) (Formatter, error) {
	return NewJSONFormatter(), nil
}
