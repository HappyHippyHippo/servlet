package log

// Formatter interface defines the methods of a logging formatter instance
// responsable to parse a logging request into the output string.
type Formatter interface {
	Format(level Level, message string, fields F) string
}
