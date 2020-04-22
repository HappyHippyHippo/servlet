package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/happyhippyhippo/servlet"
)

// LogResponseReader @TODO
type LogResponseReader interface {
	Get(context servlet.Context) map[string]interface{}
}

type logResposeReader struct{}

// NewLogResponseReader @TODO
func NewLogResponseReader() LogResponseReader {
	return &logResposeReader{}
}

// Get @TODO
func (r logResposeReader) Get(context servlet.Context) map[string]interface{} {
	response := context.(*gin.Context).Writer.(LogResponseWriter)

	return map[string]interface{}{
		"status":  response.Status(),
		"headers": r.headers(response),
		"body":    string(response.Body()),
		"time":    time.Now().Format("2006-01-02T15:04:05.000-0700"),
	}
}

func (logResposeReader) headers(response gin.ResponseWriter) map[string][]string {
	headers := map[string][]string{}
	for index, header := range response.Header() {
		headers[index] = header
	}
	return headers
}
