package middleware

import (
	"os"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/happyhippyhippo/servlet/log"
)

func Test_NewLogMiddlewareParameters(t *testing.T) {
	logChannel := "channel"
	logLevel := log.FATAL
	logRequestMessage := "request message"
	logResponseMessage := "response message"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	next := func(c *gin.Context) {}
	logger := NewMockLogger(ctrl)

	t.Run("creates a new parameters", func(t *testing.T) {
		parameters := NewLogMiddlewareParameters(next, logger)

		if reflect.ValueOf(parameters.Next).Pointer() != reflect.ValueOf(next).Pointer() {
			t.Errorf("didn't stored the next handler")
		} else if value := parameters.Logger; value != logger {
			t.Errorf("stored (%v) logger reference", value)
		} else if value := parameters.LogChannel; value != LogChannel {
			t.Errorf("stored (%v) log channel", value)
		} else if value := parameters.LogLevel; value != LogLevel {
			t.Errorf("stored (%v) log level", value)
		} else if value := parameters.LogRequestMessage; value != LogRequestMessage {
			t.Errorf("stored (%v) log request message", value)
		} else if value := parameters.LogResponseMessage; value != LogResponseMessage {
			t.Errorf("stored (%v) log response message", value)
		}
	})

	t.Run("creates a new parameters with the env log channel", func(t *testing.T) {
		os.Setenv(EnvLogChannel, logChannel)
		defer os.Setenv(EnvLogChannel, "")

		parameters := NewLogMiddlewareParameters(next, logger)
		if value := parameters.LogChannel; value != logChannel {
			t.Errorf("stored (%v) log channel", value)
		}
	})

	t.Run("creates a new parameters with a valid env log level", func(t *testing.T) {
		os.Setenv(EnvLogLevel, log.LevelNameMap[logLevel])
		defer os.Setenv(EnvLogLevel, "")

		parameters := NewLogMiddlewareParameters(next, logger)
		if value := parameters.LogLevel; value != logLevel {
			t.Errorf("stored (%v) log level", value)
		}
	})

	t.Run("error on new parameters with a invalid env log level", func(t *testing.T) {
		os.Setenv(EnvLogLevel, "invalid")
		defer os.Setenv(EnvLogLevel, "")

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		NewLogMiddlewareParameters(next, logger)
	})

	t.Run("creates a new parameters with a valid env log request message", func(t *testing.T) {
		os.Setenv(EnvLogRequestMessage, logRequestMessage)
		defer os.Setenv(EnvLogRequestMessage, "")

		parameters := NewLogMiddlewareParameters(next, logger)
		if value := parameters.LogRequestMessage; value != logRequestMessage {
			t.Errorf("stored (%v) log request message", value)
		}
	})

	t.Run("creates a new parameters with a valid env log response message", func(t *testing.T) {
		os.Setenv(EnvLogResponseMessage, logResponseMessage)
		defer os.Setenv(EnvLogResponseMessage, "")

		parameters := NewLogMiddlewareParameters(next, logger)
		if value := parameters.LogResponseMessage; value != logResponseMessage {
			t.Errorf("stored (%v) log response message", value)
		}
	})
}
