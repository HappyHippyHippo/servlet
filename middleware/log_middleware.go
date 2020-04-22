package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/happyhippyhippo/servlet"
	"github.com/happyhippyhippo/servlet/log"
)

type logMiddleware struct {
	params LogMiddlewareParameters
}

// NewLogMiddleware will instantiate a new middleware that will emit logging
// signals on a request event and on a response event.
func NewLogMiddleware(params LogMiddlewareParameters) Middleware {
	return &logMiddleware{
		params: params,
	}
}

// Execute will execute the process of logging with the related forwarding
// to the registed next handler func given on the constructor function.
func (m logMiddleware) Execute(context servlet.Context) {
	gcontext := context.(*gin.Context)
	gcontext.Writer, _ = NewLogResponseWriter(gcontext.Writer)

	request := m.params.RequestReader.Get(context)
	m.params.Logger.Signal(m.params.LogChannel, m.params.LogLevel, m.params.LogRequestMessage, log.F{"request": request})

	m.params.Next(gcontext)

	response := m.params.ResponseReader.Get(context)
	m.params.Logger.Signal(m.params.LogChannel, m.params.LogLevel, m.params.LogResponseMessage, log.F{"request": request, "response": response})
}
