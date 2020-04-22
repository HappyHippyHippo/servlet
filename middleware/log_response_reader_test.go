package middleware

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func Test_NewLogResponseReader(t *testing.T) {
	t.Run("creates a new log response request", func(t *testing.T) {
		if NewLogResponseReader() == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_LogResponseReader_Get(t *testing.T) {
	status := 200
	headers := map[string][]string{"header1": {"value1"}, "header2": {"value2"}}
	jsonBody := map[string]interface{}{"field": "value"}
	rawBody, _ := json.Marshal(jsonBody)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := NewMockLogResponseWriter(ctrl)
	writer.EXPECT().Body().Return(rawBody).Times(1)
	writer.EXPECT().Status().Return(status).Times(1)
	writer.EXPECT().Header().Return(headers).Times(1)

	context := &gin.Context{}
	context.Writer = writer

	reader := NewLogResponseReader()
	data := reader.Get(context)

	t.Run("retrieve the response status", func(t *testing.T) {
		if value := data["status"]; value != status {
			t.Errorf("stored the (%s) status value", value)
		}
	})
	t.Run("retrieve the response headers", func(t *testing.T) {
		if value := data["headers"]; !reflect.DeepEqual(value, headers) {
			t.Errorf("stored the (%v) headers", value)
		}
	})

	t.Run("retrieve the response raw body", func(t *testing.T) {
		if value := data["body"].(map[string]interface{})["raw"]; !reflect.DeepEqual(value, string(rawBody)) {
			t.Errorf("stored the (%v) raw body", value)
		}
	})

	t.Run("retrieve the response json body", func(t *testing.T) {
		if value := data["body"].(map[string]interface{})["json"]; !reflect.DeepEqual(value, jsonBody) {
			t.Errorf("stored the (%v) json body", value)
		}
	})
}
