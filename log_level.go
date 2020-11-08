package servlet

// LogLevel identifies a value type that describes a logging level.
type LogLevel int

const (
	// FATAL defines a fatal logging level.
	FATAL LogLevel = 1 + iota
	// ERROR defines a error logging level.
	ERROR
	// WARNING defines a warning logging level.
	WARNING
	// NOTICE defines a notice logging level.
	NOTICE
	// INFO defines a info logging level.
	INFO
	// DEBUG defines a debug logging level.
	DEBUG
)

// LogLevelMap defines a relation between a human-readable string
// and a code level identifier of a logging level.
var LogLevelMap = map[string]LogLevel{
	"fatal":   FATAL,
	"error":   ERROR,
	"warning": WARNING,
	"notice":  NOTICE,
	"info":    INFO,
	"debug":   DEBUG,
}

// LogLevelNameMap defines a relation between a code level identifier of a
// logging level and human-readable string representation of that level.
var LogLevelNameMap = map[LogLevel]string{
	FATAL:   "fatal",
	ERROR:   "error",
	WARNING: "warning",
	NOTICE:  "notice",
	INFO:    "info",
	DEBUG:   "debug",
}
