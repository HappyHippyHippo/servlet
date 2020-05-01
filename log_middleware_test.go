package servlet

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
)

/// ---------------------------------------------------------------------------
/// LogMiddlewareBasicRequestReader
/// ---------------------------------------------------------------------------

func Test_NewLogMiddlewareBasicRequestReader(t *testing.T) {
	t.Run("new log request request", func(t *testing.T) {
		if NewLogMiddlewareBasicRequestReader() == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_LogMiddlewareBasicRequestReader_Get(t *testing.T) {
	method := "method"
	uri := "/resource"
	reqUrl, _ := url.Parse("http://domain" + uri)
	headers := map[string][]string{"header1": {"value1"}, "header2": {"value2"}}
	jsonBody := map[string]interface{}{"field": "value"}
	rawBody, _ := json.Marshal(jsonBody)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := NewMockReader(ctrl)
	gomock.InOrder(
		body.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) { copy(p, rawBody); return len(rawBody), nil }),
		body.EXPECT().Read(gomock.Any()).Return(0, io.EOF),
	)

	context := &gin.Context{}
	context.Request = &http.Request{}
	context.Request.Method = method
	context.Request.URL = reqUrl
	context.Request.Header = headers
	context.Request.Body = body

	reader := NewLogMiddlewareBasicRequestReader()
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

	t.Run("retrieve the request body", func(t *testing.T) {
		if value := data["body"]; !reflect.DeepEqual(value, string(rawBody)) {
			t.Errorf("stored the (%v) body", value)
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareJsonRequestReader
/// ---------------------------------------------------------------------------

func Test_NewLogMiddlewareJsonRequestReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockLogMiddlewareRequestReader(ctrl)

	t.Run("nil reader", func(t *testing.T) {
		if reader, err := NewLogMiddlewareJsonRequestReader(nil, nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'reader' argument" {
			t.Errorf("returned the (%v) error", err)
		} else if reader != nil {
			t.Error("returned a valid reference")
		}
	})

	t.Run("new decorator", func(t *testing.T) {
		if reader, err := NewLogMiddlewareJsonRequestReader(reader, nil); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Error("didn't returned a valid reference")
		}
	})

	t.Run("new decorator with model", func(t *testing.T) {
		model := struct{ data string }{data: "bing"}

		if reader, err := NewLogMiddlewareJsonRequestReader(reader, model); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Error("didn't returned a valid reference")
		} else if check := reader.model; check != model {
			t.Errorf("stored the (%v) model", check)
		}
	})
}

func Test_LogMiddlewareJsonRequestReader_Get(t *testing.T) {
	t.Run("non-json body does not add decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "{"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockGinContext(ctrl)
		reader := NewMockLogMiddlewareRequestReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		jsonReader, _ := NewLogMiddlewareJsonRequestReader(reader, nil)
		result := jsonReader.Get(context)

		if _, ok := result["bodyJson"]; ok {
			t.Error("added the bodyJson field")
		}
	})

	t.Run("json body adds decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "{\"field\":\"value\"}"}
		expected := map[string]interface{}{"field": "value"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockGinContext(ctrl)
		reader := NewMockLogMiddlewareRequestReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		jsonReader, _ := NewLogMiddlewareJsonRequestReader(reader, nil)
		result := jsonReader.Get(context)

		if body, ok := result["bodyJson"]; !ok {
			t.Error("didn't added the bodyJson field")
		} else if !reflect.DeepEqual(body, expected) {
			t.Errorf("added the (%v) json content)", body)
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareXmlRequestReader
/// ---------------------------------------------------------------------------

func Test_NewLogMiddlewareXmlRequestReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockLogMiddlewareRequestReader(ctrl)

	t.Run("error if reader is nil", func(t *testing.T) {
		if reader, err := NewLogMiddlewareXmlRequestReader(nil, nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'reader' argument" {
			t.Errorf("returned the (%v) error", err)
		} else if reader != nil {
			t.Error("returned a valid reference")
		}
	})

	t.Run("new decorator", func(t *testing.T) {
		if reader, err := NewLogMiddlewareXmlRequestReader(reader, nil); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Error("didn't return a valid reference")
		}
	})

	t.Run("new decorator with model", func(t *testing.T) {
		model := struct{ data string }{data: "bing"}

		if reader, err := NewLogMiddlewareXmlRequestReader(reader, model); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Error("didn't returned a valid reference")
		} else if check := reader.model; check != model {
			t.Errorf("stored the (%v) model", check)
		}
	})
}

func Test_LogMiddlewareXmlRequestReader_Get(t *testing.T) {
	model := struct {
		XMLName xml.Name `xml:"message"`
		Field   string   `xml:"field"`
	}{}

	t.Run("non-xml body does not add decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "{"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockGinContext(ctrl)
		reader := NewMockLogMiddlewareRequestReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		xmlReader, _ := NewLogMiddlewareXmlRequestReader(reader, &model)
		result := xmlReader.Get(context)

		if _, ok := result["bodyXml"]; ok {
			t.Errorf("added the bodyXml field")
		}
	})

	t.Run("xml body adds decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "\u003cmessage\u003e\u003cfield\u003evalue\u003c/field\u003e\u003c/message\u003e"}
		expected := struct {
			XMLName xml.Name `xml:"message"`
			Field   string   `xml:"field"`
		}{XMLName: xml.Name{Local: "message"}, Field: "value"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockGinContext(ctrl)
		reader := NewMockLogMiddlewareRequestReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		xmlReader, _ := NewLogMiddlewareXmlRequestReader(reader, &model)
		result := xmlReader.Get(context)

		if body, ok := result["bodyXml"]; !ok {
			t.Errorf("didn't added the bodyXml field")
		} else if !reflect.DeepEqual(body, &expected) {
			t.Errorf("added the (%v) xml content)", body)
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareResponseWriter
/// ---------------------------------------------------------------------------

func Test_NewLogMiddlewareResponseWriter(t *testing.T) {
	t.Run("error when missing writer", func(t *testing.T) {
		if value, err := NewLogMiddlewareResponseWriter(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'writer' argument" {
			t.Errorf("returned the (%v) error", err)
		} else if value != nil {
			t.Error("returned a valid reference")
		}
	})

	t.Run("new log response writer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		writer := NewMockResponseWriter(ctrl)

		if value, err := NewLogMiddlewareResponseWriter(writer); err != nil {
			t.Errorf("return the (%v) error", err)
		} else if value == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_LogMiddlewareResponseWriter_Write(t *testing.T) {
	b := []byte{12, 34, 56}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ginWriter := NewMockResponseWriter(ctrl)
	ginWriter.EXPECT().Write(b).Times(1)

	writer := &LogMiddlewareResponseWriter{body: &bytes.Buffer{}, ResponseWriter: ginWriter}

	t.Run("write to buffer and underlying writer", func(t *testing.T) {
		_, _ = writer.Write(b)
		if !reflect.DeepEqual(writer.body.Bytes(), b) {
			t.Errorf("written (%v) bytes on buffer", writer.body)
		}
	})
}

func Test_LogResponseWriter_Body(t *testing.T) {
	b := []byte{12, 34, 56}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ginWriter := NewMockResponseWriter(ctrl)
	writer := &LogMiddlewareResponseWriter{body: bytes.NewBuffer(b), ResponseWriter: ginWriter}

	t.Run("write to buffer and underlying writer", func(t *testing.T) {
		if !reflect.DeepEqual(writer.Body(), b) {
			t.Errorf("written (%v) bytes on buffer", writer.body)
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareBasicResponseReader
/// ---------------------------------------------------------------------------

func Test_NewLogMiddlewareBasicResponseReader(t *testing.T) {
	t.Run("new log response request", func(t *testing.T) {
		if NewLogMiddlewareBasicResponseReader() == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_LogMiddlewareBasicResponseReader_Get(t *testing.T) {
	status := 200
	headers := map[string][]string{"header1": {"value1"}, "header2": {"value2"}}
	jsonBody := map[string]interface{}{"field": "value"}
	rawBody, _ := json.Marshal(jsonBody)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ginWriter := NewMockResponseWriter(ctrl)
	ginWriter.EXPECT().Status().Return(status).Times(1)
	ginWriter.EXPECT().Header().Return(headers).Times(1)

	writer, _ := NewLogMiddlewareResponseWriter(ginWriter)
	writer.body.Write(rawBody)

	context := &gin.Context{}
	context.Writer = writer

	reader := NewLogMiddlewareBasicResponseReader()
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

	t.Run("retrieve the response body", func(t *testing.T) {
		if value := data["body"]; !reflect.DeepEqual(value, string(rawBody)) {
			t.Errorf("stored the (%v) body", value)
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareJsonResponseReader
/// ---------------------------------------------------------------------------

func Test_NewLogMiddlewareJsonResponseReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockLogMiddlewareResponseReader(ctrl)

	t.Run("nil reader", func(t *testing.T) {
		if reader, err := NewLogMiddlewareJsonResponseReader(nil, nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'reader' argument" {
			t.Errorf("returned the (%v) error", err)
		} else if reader != nil {
			t.Error("returned a valid reference")
		}
	})

	t.Run("new decorator", func(t *testing.T) {
		if reader, err := NewLogMiddlewareJsonResponseReader(reader, nil); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Error("didn't returned a valid reference")
		}
	})

	t.Run("new decorator with model", func(t *testing.T) {
		model := struct{ data string }{data: "bing"}

		if reader, err := NewLogMiddlewareJsonResponseReader(reader, model); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Error("didn't returned a valid reference")
		} else if check := reader.model; check != model {
			t.Errorf("stored the (%v) model", check)
		}
	})
}

func Test_LogMiddlewareJsonResponseReader_Get(t *testing.T) {
	t.Run("non-json body does not add decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "{"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockGinContext(ctrl)
		reader := NewMockLogMiddlewareResponseReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		jsonReader, _ := NewLogMiddlewareJsonResponseReader(reader, nil)
		result := jsonReader.Get(context)

		if _, ok := result["bodyJson"]; ok {
			t.Error("added the bodyJson field")
		}
	})

	t.Run("json body adds decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "{\"field\":\"value\"}"}
		expected := map[string]interface{}{"field": "value"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockGinContext(ctrl)
		reader := NewMockLogMiddlewareResponseReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		jsonReader, _ := NewLogMiddlewareJsonResponseReader(reader, nil)
		result := jsonReader.Get(context)

		if body, ok := result["bodyJson"]; !ok {
			t.Error("didn't added the bodyJson field")
		} else if !reflect.DeepEqual(body, expected) {
			t.Errorf("added the (%v) json content)", body)
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareXmlResponseReader
/// ---------------------------------------------------------------------------

func Test_NewLogMiddlewareXmlResponseReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockLogMiddlewareResponseReader(ctrl)

	t.Run("nil reader", func(t *testing.T) {
		if reader, err := NewLogMiddlewareXmlResponseReader(nil, nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'reader' argument" {
			t.Errorf("returned the (%v) error", err)
		} else if reader != nil {
			t.Error("returned a valid reference")
		}
	})

	t.Run("new decorator", func(t *testing.T) {
		if reader, err := NewLogMiddlewareXmlResponseReader(reader, nil); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Error("didn't returned a valid reference")
		}
	})

	t.Run("new decorator with model", func(t *testing.T) {
		model := struct{ data string }{data: "bing"}

		if reader, err := NewLogMiddlewareXmlResponseReader(reader, model); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Error("didn't return a valid reference")
		} else if check := reader.model; check != model {
			t.Errorf("stored the (%v) model", check)
		}
	})
}

func Test_LogMiddlewareXmlResponseReader_Get(t *testing.T) {
	model := struct {
		XMLName xml.Name `xml:"message"`
		Field   string   `xml:"field"`
	}{}

	t.Run("non-xml body does not add decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "{"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockGinContext(ctrl)
		reader := NewMockLogMiddlewareResponseReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		xmlReader, _ := NewLogMiddlewareXmlResponseReader(reader, &model)
		result := xmlReader.Get(context)

		if _, ok := result["bodyXml"]; ok {
			t.Error("added the bodyXml field")
		}
	})

	t.Run("xml body adds decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "\u003cmessage\u003e\u003cfield\u003evalue\u003c/field\u003e\u003c/message\u003e"}
		expected := struct {
			XMLName xml.Name `xml:"message"`
			Field   string   `xml:"field"`
		}{XMLName: xml.Name{Local: "message"}, Field: "value"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockGinContext(ctrl)
		reader := NewMockLogMiddlewareResponseReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		xmlReader, _ := NewLogMiddlewareXmlResponseReader(reader, &model)
		result := xmlReader.Get(context)

		if body, ok := result["bodyXml"]; !ok {
			t.Error("didn't added the bodyXml field")
		} else if !reflect.DeepEqual(body, &expected) {
			t.Errorf("added the (%v) xml content)", body)
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareParams
/// ---------------------------------------------------------------------------

func Test_NewLogMiddlewareParams(t *testing.T) {
	next := func(c *gin.Context) {}
	logger := NewLog()

	t.Run("nil next", func(t *testing.T) {
		if parameters, err := NewLogMiddlewareParams(nil, logger); parameters != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'next' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil logger", func(t *testing.T) {
		if parameters, err := NewLogMiddlewareParams(next, nil); parameters != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'logger' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new parameters", func(t *testing.T) {
		if parameters, err := NewLogMiddlewareParams(next, logger); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reflect.ValueOf(parameters.Next).Pointer() != reflect.ValueOf(next).Pointer() {
			t.Error("didn't stored the next handler")
		} else if value := parameters.Logger; value != logger {
			t.Errorf("stored (%v) logger reference", value)
		} else if value := parameters.LogChannel; value != LogMiddlewareChannel {
			t.Errorf("stored (%v) logging channel", value)
		} else if value := parameters.LogLevel; value != LogMiddlewareLevel {
			t.Errorf("stored (%v) logging level", value)
		} else if value := parameters.LogRequestMessage; value != LogMiddlewareRequestMessage {
			t.Errorf("stored (%v) logging request message", value)
		} else if value := parameters.LogResponseMessage; value != LogMiddlewareResponseMessage {
			t.Errorf("stored (%v) logging response message", value)
		}
	})

	t.Run("new parameters with the env log channel", func(t *testing.T) {
		logChannel := "channel"
		_ = os.Setenv(EnvLogMiddlewareChannel, logChannel)
		defer func() { _ = os.Setenv(EnvLogMiddlewareChannel, "") }()

		parameters, _ := NewLogMiddlewareParams(next, logger)
		if value := parameters.LogChannel; value != logChannel {
			t.Errorf("stored (%v) log channel", value)
		}
	})

	t.Run("new parameters with a valid env log level", func(t *testing.T) {
		logLevel := FATAL
		_ = os.Setenv(EnvLogMiddlewareLevel, LogLevelNameMap[logLevel])
		defer func() { _ = os.Setenv(EnvLogMiddlewareLevel, "") }()

		parameters, _ := NewLogMiddlewareParams(next, logger)
		if value := parameters.LogLevel; value != logLevel {
			t.Errorf("stored (%v) log level", value)
		}
	})

	t.Run("error on new parameters with a invalid env log level", func(t *testing.T) {
		_ = os.Setenv(EnvLogMiddlewareLevel, "invalid")
		defer func() { _ = os.Setenv(EnvLogMiddlewareLevel, "") }()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		_, _ = NewLogMiddlewareParams(next, logger)
	})

	t.Run("new parameters with a valid env log request message", func(t *testing.T) {
		logRequestMessage := "request message"
		_ = os.Setenv(EnvLogMiddlewareRequestMessage, logRequestMessage)
		defer func() { _ = os.Setenv(EnvLogMiddlewareRequestMessage, "") }()

		parameters, _ := NewLogMiddlewareParams(next, logger)
		if value := parameters.LogRequestMessage; value != logRequestMessage {
			t.Errorf("stored (%v) log request message", value)
		}
	})

	t.Run("new parameters with a valid env log response message", func(t *testing.T) {
		logResponseMessage := "response message"
		_ = os.Setenv(EnvLogMiddlewareResponseMessage, logResponseMessage)
		defer func() { _ = os.Setenv(EnvLogMiddlewareResponseMessage, "") }()

		parameters, _ := NewLogMiddlewareParams(next, logger)
		if value := parameters.LogResponseMessage; value != logResponseMessage {
			t.Errorf("stored (%v) log response message", value)
		}
	})
}

/// ---------------------------------------------------------------------------
/// LogMiddleware
/// ---------------------------------------------------------------------------

func Test_NewLogMiddleware(t *testing.T) {
	t.Run("nil parameters", func(t *testing.T) {
		if middleware, err := NewLogMiddleware(nil); middleware != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'parameters' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := NewMockResponseWriter(ctrl)
	context := &gin.Context{}
	context.Writer = writer

	callCount := 0
	var next gin.HandlerFunc = func(*gin.Context) { callCount = callCount + 1 }

	request := map[string]interface{}{"type": "request"}
	requestReader := NewMockLogMiddlewareRequestReader(ctrl)
	requestReader.EXPECT().Get(context).Return(request).Times(1)

	response := map[string]interface{}{"type": "response"}
	responseReader := NewMockLogMiddlewareResponseReader(ctrl)
	responseReader.EXPECT().Get(context).Return(response).Times(1)

	logStream := NewMockLogStream(ctrl)
	gomock.InOrder(
		logStream.EXPECT().Signal(LogMiddlewareChannel, LogMiddlewareLevel, LogMiddlewareRequestMessage, map[string]interface{}{"request": request}),
		logStream.EXPECT().Signal(LogMiddlewareChannel, LogMiddlewareLevel, LogMiddlewareResponseMessage, map[string]interface{}{"request": request, "response": response}),
	)
	logger := NewLog()
	_ = logger.AddStream("id", logStream)

	parameters, _ := NewLogMiddlewareParams(next, logger)
	parameters.Logger = logger
	parameters.Next = next
	parameters.RequestReader = requestReader
	parameters.ResponseReader = responseReader

	mw, _ := NewLogMiddleware(parameters)
	mw(context)

	t.Run("call next handler", func(t *testing.T) {
		if callCount != 1 {
			t.Errorf("didn't called the next handler")
		}
	})
}
