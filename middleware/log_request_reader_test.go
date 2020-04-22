package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func Test_NewLogRequestReader(t *testing.T) {
	t.Run("creates a new log request request", func(t *testing.T) {
		if NewLogRequestReader() == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_LogRequestReader_Get(t *testing.T) {
	method := "method"
	uri := "/resource"
	url, _ := url.Parse("http://domain" + uri)
	headers := map[string][]string{"header1": {"value1"}, "header2": {"value2"}}
	jsonBody := map[string]interface{}{"field": "value"}
	rawBody, _ := json.Marshal(jsonBody)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := NewMockReadCloser(ctrl)
	gomock.InOrder(
		body.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) { copy(p, rawBody); return len(rawBody), nil }),
		body.EXPECT().Read(gomock.Any()).Return(0, io.EOF),
	)

	context := &gin.Context{}
	context.Request = &http.Request{}
	context.Request.Method = method
	context.Request.URL = url
	context.Request.Header = headers
	context.Request.Body = body

	reader := NewLogRequestReader()
	data := reader.Get(context)

	t.Run("retrieve the request method", func(t *testing.T) {
		if value := data["method"]; value != method {
			t.Errorf("stored the (%s) method value", value)
		}
	})

	t.Run("retrieve the request URI", func(t *testing.T) {
		if value := data["uri"]; value != uri {
			t.Errorf("stored the (%s) uri value", value)
		}
	})

	t.Run("retrieve the request headers", func(t *testing.T) {
		if value := data["headers"]; !reflect.DeepEqual(value, headers) {
			t.Errorf("stored the (%v) headers", value)
		}
	})

	t.Run("retrieve the request raw body", func(t *testing.T) {
		if value := data["body"].(map[string]interface{})["raw"]; !reflect.DeepEqual(value, string(rawBody)) {
			t.Errorf("stored the (%v) raw body", value)
		}
	})

	t.Run("retrieve the request json body", func(t *testing.T) {
		if value := data["body"].(map[string]interface{})["json"]; !reflect.DeepEqual(value, jsonBody) {
			t.Errorf("stored the (%v) json body", value)
		}
	})
}
