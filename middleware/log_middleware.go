package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/happyhippyhippo/servlet"
	"github.com/happyhippyhippo/servlet/log"
)

type logMiddleware struct {
	params LogMiddlewareParameters
}

// NewLogMiddleware @TODO
func NewLogMiddleware(params LogMiddlewareParameters) Middleware {
	return &logMiddleware{
		params: params,
	}
}

// Execute @TODO
func (m logMiddleware) Execute(context servlet.Context) {
	gcontext := context.(*gin.Context)
	gcontext.Writer, _ = NewLogResponseWriter(gcontext.Writer)

	request := m.params.RequestReader.Get(context)
	m.params.Logger.Signal(m.params.LogChannel, m.params.LogLevel, m.params.LogRequestMessage, log.F{"request": request})

	m.params.Next(gcontext)

	response := m.params.ResponseReader.Get(context)
	m.params.Logger.Signal(m.params.LogChannel, m.params.LogLevel, m.params.LogResponseMessage, log.F{"request": request, "response": response})
}
