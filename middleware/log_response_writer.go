package middleware

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
)

// LogResponseWriter @TODO
type LogResponseWriter interface {
	gin.ResponseWriter
	Body() []byte
}

type logResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// NewLogResponseWriter @TODO
func NewLogResponseWriter(writer gin.ResponseWriter) (LogResponseWriter, error) {
	if writer == nil {
		return nil, fmt.Errorf("Invalid nil 'writer' argument")
	}

	return &logResponseWriter{
		body:           &bytes.Buffer{},
		ResponseWriter: writer,
	}, nil
}

// Write @TODO
func (w logResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Body @TODO
func (w logResponseWriter) Body() []byte {
	return w.body.Bytes()
}
