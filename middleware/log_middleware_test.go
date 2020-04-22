package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/happyhippyhippo/servlet/log"
)

func Test_NewLogMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var next gin.HandlerFunc = func(*gin.Context) {}
	logger := NewMockLogger(ctrl)

	parameters := NewLogMiddlewareParameters(next, logger)

	t.Run("creates a new log middleware", func(t *testing.T) {
		if mw := NewLogMiddleware(parameters).(*logMiddleware); mw == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_LogMiddleware_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := NewMockGinResponseWriter(ctrl)
	context := &gin.Context{}
	context.Writer = writer

	callCount := 0
	var next gin.HandlerFunc = func(*gin.Context) { callCount = callCount + 1 }

	request := map[string]interface{}{"type": "request"}
	reqReader := NewMockLogRequestReader(ctrl)
	reqReader.EXPECT().Get(context).Return(request).Times(1)

	response := map[string]interface{}{"type": "response"}
	resReader := NewMockLogResponseReader(ctrl)
	resReader.EXPECT().Get(context).Return(response).Times(1)

	logger := NewMockLogger(ctrl)
	gomock.InOrder(
		logger.EXPECT().Signal(LogChannel, LogLevel, LogRequestMessage, log.F{"request": request}),
		logger.EXPECT().Signal(LogChannel, LogLevel, LogResponseMessage, log.F{"request": request, "response": response}),
	)

	parameters := NewLogMiddlewareParameters(next, logger)
	parameters.Logger = logger
	parameters.Next = next
	parameters.ReqReader = reqReader
	parameters.ResReader = resReader

	mw := NewLogMiddleware(parameters)

	mw.Execute(context)

	t.Run("call next handler", func(t *testing.T) {
		if callCount != 1 {
			t.Errorf("didn't called the next handler")
		}
	})
}
