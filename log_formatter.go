package servlet

// LogFormatter interface defines the methods of a logging formatter instance
// responsible to parse a logging request into the output string.
type LogFormatter interface {
	Format(level LogLevel, message string, context map[string]interface{}) string
}
