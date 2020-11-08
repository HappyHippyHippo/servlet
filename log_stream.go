package servlet

// LogStream interface defines the interaction methods with a logging stream.
type LogStream interface {
	Close() error
	Signal(channel string, level LogLevel, message string, context map[string]interface{}) error
	Broadcast(level LogLevel, message string, context map[string]interface{}) error
	HasChannel(channel string) bool
	ListChannels() []string
	AddChannel(channel string)
	RemoveChannel(channel string)
	Level() LogLevel
}
