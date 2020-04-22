package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/happyhippyhippo/servlet/log"
)

const (
	// LogChannel @TODO
	LogChannel = "conn"

	// EnvLogChannel @TODO
	EnvLogChannel = "SERVLET_LOG_MIDDLEWARE_CHANNEL"

	// LogLevel @TODO
	LogLevel = log.INFO

	// EnvLogLevel @TODO
	EnvLogLevel = "SERVLET_LOG_MIDDLEWARE_LEVEL"

	// LogRequestMessage @TODO
	LogRequestMessage = "Request"

	// EnvLogRequestMessage @TODO
	EnvLogRequestMessage = "SERVLET_LOG_MIDDLEWARE_REQUEST_MESSAGE"

	// LogResponseMessage @TODO
	LogResponseMessage = "Response"

	// EnvLogResponseMessage @TODO
	EnvLogResponseMessage = "SERVLET_LOG_MIDDLEWARE_RESPONSE_MESSAGE"
)

// LogMiddlewareParameters @type
type LogMiddlewareParameters struct {
	ReqReader          LogRequestReader
	ResReader          LogResponseReader
	Next               gin.HandlerFunc
	Logger             log.Logger
	LogChannel         string
	LogLevel           log.Level
	LogRequestMessage  string
	LogResponseMessage string
}

// NewLogMiddlewareParameters @TODO
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
		ReqReader:          NewLogRequestReader(),
		ResReader:          NewLogResponseReader(),
		Next:               next,
		Logger:             logger,
		LogChannel:         logChannel,
		LogLevel:           logLevel,
		LogRequestMessage:  logRequestMessage,
		LogResponseMessage: logResponseMessage,
	}
}
