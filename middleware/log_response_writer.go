package middleware

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
)

// LogResponseWriter defines a response writer proxy instance used to store the
// response content so it can be used for response composition on the
// logging process.
type LogResponseWriter interface {
	gin.ResponseWriter
	Body() []byte
}

type logResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// NewLogResponseWriter intantiate a new response writer proxy.
func NewLogResponseWriter(writer gin.ResponseWriter) (LogResponseWriter, error) {
	if writer == nil {
		return nil, fmt.Errorf("Invalid nil 'writer' argument")
	}

	return &logResponseWriter{
		body:           &bytes.Buffer{},
		ResponseWriter: writer,
	}, nil
}

// Write executes the writing the desired bytes into the underlying writer
// and storing them in the internal buffer.
func (w logResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Body will retrieve the stored bytes given on the previous calls
// to the Write method.
func (w logResponseWriter) Body() []byte {
	return w.body.Bytes()
}
