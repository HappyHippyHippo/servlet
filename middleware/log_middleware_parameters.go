package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/happyhippyhippo/servlet/log"
)

const (
	// LogChannel defines the channel id to be used when the log middleware
	// sends the logging signal to the logger instance.
	LogChannel = "conn"

	// EnvLogChannel defines the name of the environment
	// variable to be checked for a overriding value for the channel id used
	// on the logging signal call.
	EnvLogChannel = "SERVLET_LOG_MIDDLEWARE_CHANNEL"

	// LogLevel defines the logging level to be used when the log middleware
	// sends the logging signal to the logger instance.
	LogLevel = log.INFO

	// EnvLogLevel defines the name of the environment
	// variable to be checked for a overriding value for the logging level
	// used on the logging signal call.
	EnvLogLevel = "SERVLET_LOG_MIDDLEWARE_LEVEL"

	// LogRequestMessage defines the request event logging message to be used
	// when the log middleware sends the logging signal to the logger instance.
	LogRequestMessage = "Request"

	// EnvLogRequestMessage defines the name of the environment
	// variable to be checked for a overriding value for the request event
	// logging message used on the logging signal call
	EnvLogRequestMessage = "SERVLET_LOG_MIDDLEWARE_REQUEST_MESSAGE"

	// LogResponseMessage defines the response event logging message to be used
	// when the log middleware sends the logging signal to the logger instance.
	LogResponseMessage = "Response"

	// EnvLogResponseMessage defines the name of the environment
	// variable to be checked for a overriding value for the response event
	// logging message used on the logging signal call
	EnvLogResponseMessage = "SERVLET_LOG_MIDDLEWARE_RESPONSE_MESSAGE"
)

// LogMiddlewareParameters defines the storing structure of the parameters
// used to configure the logging middleware.
type LogMiddlewareParameters struct {
	RequestReader      LogRequestReader
	ResponseReader     LogResponseReader
	Next               gin.HandlerFunc
	Logger             log.Logger
	LogChannel         string
	LogLevel           log.Level
	LogRequestMessage  string
	LogResponseMessage string
}

// NewLogMiddlewareParameters will instantiate a new log middleware parameters
// instance used to configure a log middleware. If environment variables have
// been set for the log environment, the returned parameters structure will
// reflect those values.
func NewLogMiddlewareParameters(next gin.HandlerFunc, logger log.Logger) LogMiddlewareParameters {
	logChannel := LogChannel
	if env := os.Getenv(EnvLogChannel); env != "" {
		logChannel = env
	}

	logLevel := LogLevel
	if env := os.Getenv(EnvLogLevel); env != "" {
		env = strings.ToLower(env)
		if l, ok := log.LevelMap[env]; !ok {
			panic(fmt.Errorf("Unrecognized logger level : %s", env))
		} else {
			logLevel = l
		}
	}

	logRequestMessage := LogRequestMessage
	if env := os.Getenv(EnvLogRequestMessage); env != "" {
		logRequestMessage = env
	}

	logResponseMessage := LogResponseMessage
	if env := os.Getenv(EnvLogResponseMessage); env != "" {
		logResponseMessage = env
	}

	return LogMiddlewareParameters{
		RequestReader:      NewLogRequestReader(),
		ResponseReader:     NewLogResponseReader(),
		Next:               next,
		Logger:             logger,
		LogChannel:         logChannel,
		LogLevel:           logLevel,
		LogRequestMessage:  logRequestMessage,
		LogResponseMessage: logResponseMessage,
	}
}
