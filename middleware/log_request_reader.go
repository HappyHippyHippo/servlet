package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/happyhippyhippo/servlet"
)

// LogRequestReader @TODO
type LogRequestReader interface {
	Get(context servlet.Context) map[string]interface{}
}

type logRequestReader struct{}

// NewLogRequestReader @TODO
func NewLogRequestReader() LogRequestReader {
	return &logRequestReader{}
}

// Get @TODO
func (r logRequestReader) Get(context servlet.Context) map[string]interface{} {
	request := context.(*gin.Context).Request

	return map[string]interface{}{
		"headers": r.headers(request),
		"method":  request.Method,
		"uri":     request.URL.RequestURI(),
		"body":    r.body(request),
		"time":    time.Now().Format("2006-01-02T15:04:05.000-0700"),
	}
}

func (logRequestReader) headers(request *http.Request) map[string][]string {
	headers := map[string][]string{}
	for index, header := range request.Header {
		headers[index] = header
	}
	return headers
}

func (logRequestReader) body(request *http.Request) string {
	var bodyBytes []byte
	if request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(request.Body)
	}
	request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}
