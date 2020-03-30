package log

// Formatter interface defines the methods of a logging formatter instance.
type Formatter interface {
	Format(level Level, fields F, message string) string
}
