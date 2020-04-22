package middleware

import (
	"bytes"
	"encoding/json"
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

	var bytesBody []byte = r.body(request)
	var jsonBody interface{}
	json.Unmarshal(bytesBody, &jsonBody)

	return map[string]interface{}{
		"headers": r.headers(request),
		"method":  request.Method,
		"uri":     request.URL.RequestURI(),
		"body":    map[string]interface{}{"raw": string(bytesBody), "json": jsonBody},
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

func (logRequestReader) body(request *http.Request) []byte {
	var bodyBytes []byte
	if request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(request.Body)
	}
	request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}
