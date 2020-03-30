package log

// Level identifies a value type that describes a logging level.
type Level int

const (
	// FATAL defines a fatal logging level.
	FATAL Level = 1 + iota
	// ERROR defines a error logging level.
	ERROR
	// WARNING defines a warning logging level.
	WARNING
	// NOTICE defines a fatanotice logging level.
	NOTICE
	// DEBUG defines a debug logging level.
	DEBUG
)

// LevelMap defines a relation between a human-readable string
// and a logging level.
var LevelMap = map[string]Level{
	"fatal":   FATAL,
	"error":   ERROR,
	"warning": WARNING,
	"notice":  NOTICE,
	"debug":   DEBUG,
}
