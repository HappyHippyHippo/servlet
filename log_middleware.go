package servlet

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

/// ---------------------------------------------------------------------------
/// constants
/// ---------------------------------------------------------------------------

const (
	// LogMiddlewareChannel defines the channel id to be used when the log
	// middleware sends the logging signal to the logger instance.
	LogMiddlewareChannel = "conn"

	// EnvLogMiddlewareChannel defines the name of the environment
	// variable to be checked for a overriding value for the channel id used
	// on the logging signal call.
	EnvLogMiddlewareChannel = "SERVLET_LOG_MIDDLEWARE_CHANNEL"

	// LogMiddlewareLevel defines the logging level to be used when the log
	// middleware sends the logging signal to the logger instance.
	LogMiddlewareLevel = INFO

	// EnvLogMiddlewareLevel defines the name of the environment
	// variable to be checked for a overriding value for the logging level
	// used on the logging signal call.
	EnvLogMiddlewareLevel = "SERVLET_LOG_MIDDLEWARE_LEVEL"

	// LogMiddlewareRequestMessage defines the request event logging message to
	// be used when the log middleware sends the logging signal to the logger
	// instance.
	LogMiddlewareRequestMessage = "Request"

	// EnvLogMiddlewareRequestMessage defines the name of the environment
	// variable to be checked for a overriding value for the request event
	// logging message used on the logging signal call
	EnvLogMiddlewareRequestMessage = "SERVLET_LOG_MIDDLEWARE_REQUEST_MESSAGE"

	// LogMiddlewareResponseMessage defines the response event logging message
	// to be used when the log middleware sends the logging signal to the
	// logger instance.
	LogMiddlewareResponseMessage = "Response"

	// EnvLogMiddlewareResponseMessage defines the name of the environment
	// variable to be checked for a overriding value for the response event
	// logging message used on the logging signal call
	EnvLogMiddlewareResponseMessage = "SERVLET_LOG_MIDDLEWARE_RESPONSE_MESSAGE"
)

/// ---------------------------------------------------------------------------
/// LogMiddlewareRequestReader
/// ---------------------------------------------------------------------------

type LogMiddlewareRequestReader interface {
	Get(context GinContext) map[string]interface{}
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareBasicRequestReader
/// ---------------------------------------------------------------------------

// LogMiddlewareBasicRequestReader defines the interface methods of a request
// context reader used to compose the data to be sent to the logger on a
// request event.
type LogMiddlewareBasicRequestReader struct{}

// NewLogMiddlewareBasicRequestReader will instantiate a new basic request context reader.
func NewLogMiddlewareBasicRequestReader() *LogMiddlewareBasicRequestReader {
	return &LogMiddlewareBasicRequestReader{}
}

// Get process the context request and return the data to be
// signaled to the logger.
func (r LogMiddlewareBasicRequestReader) Get(context GinContext) map[string]interface{} {
	request := context.(*gin.Context).Request

	return map[string]interface{}{
		"headers": r.headers(request),
		"method":  request.Method,
		"uri":     request.URL.RequestURI(),
		"body":    r.body(request),
		"time":    time.Now().Format("2006-01-02T15:04:05.000-0700"),
	}
}

func (LogMiddlewareBasicRequestReader) headers(request *http.Request) map[string][]string {
	headers := map[string][]string{}
	for index, header := range request.Header {
		headers[index] = header
	}
	return headers
}

func (LogMiddlewareBasicRequestReader) body(request *http.Request) string {
	var bodyBytes []byte
	if request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(request.Body)
	}
	request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareJsonRequestReader
/// ---------------------------------------------------------------------------

type LogMiddlewareJsonRequestReader struct {
	reader LogMiddlewareRequestReader
	model  interface{}
}

// NewLogMiddlewareJsonRequestReader will instantiate a new request event context
// reader JSON decorator used to parse the request body as a JSON and add
// the parsed content into the logging data.
func NewLogMiddlewareJsonRequestReader(reader LogMiddlewareRequestReader, model interface{}) (*LogMiddlewareJsonRequestReader, error) {
	if reader == nil {
		return nil, fmt.Errorf("invalid nil 'reader' argument")
	}

	return &LogMiddlewareJsonRequestReader{
		reader: reader,
		model:  model,
	}, nil
}

// Get process the context request and add the extra bodyJson if the body
// content can be parsed as JSON.
func (r LogMiddlewareJsonRequestReader) Get(context GinContext) map[string]interface{} {
	data := r.reader.Get(context)

	if err := json.Unmarshal([]byte(data["body"].(string)), &r.model); err == nil {
		data["bodyJson"] = r.model
	}

	return data
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareXmlRequestReader
/// ---------------------------------------------------------------------------

type LogMiddlewareXmlRequestReader struct {
	reader LogMiddlewareRequestReader
	model  interface{}
}

// NewLogMiddlewareXmlRequestReader will instantiate a new request event context
// reader XML decorator used to parse the request body as a XML and add
// the parsed content into the logging data.
func NewLogMiddlewareXmlRequestReader(reader LogMiddlewareRequestReader, model interface{}) (*LogMiddlewareXmlRequestReader, error) {
	if reader == nil {
		return nil, fmt.Errorf("invalid nil 'reader' argument")
	}

	return &LogMiddlewareXmlRequestReader{
		reader: reader,
		model:  model,
	}, nil
}

// Get process the context request and add the extra bodyJson if the body
// content can be parsed as XML.
func (r LogMiddlewareXmlRequestReader) Get(context GinContext) map[string]interface{} {
	data := r.reader.Get(context)

	if err := xml.Unmarshal([]byte(data["body"].(string)), &r.model); err == nil {
		data["bodyXml"] = r.model
	}

	return data
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareResponseWriter
/// ---------------------------------------------------------------------------

type LogMiddlewareResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// NewLogMiddlewareResponseWriter instantiate a new response writer proxy.
func NewLogMiddlewareResponseWriter(writer gin.ResponseWriter) (*LogMiddlewareResponseWriter, error) {
	if writer == nil {
		return nil, fmt.Errorf("invalid nil 'writer' argument")
	}

	return &LogMiddlewareResponseWriter{
		body:           &bytes.Buffer{},
		ResponseWriter: writer,
	}, nil
}

// Write executes the writing the desired bytes into the underlying writer
// and storing them in the internal buffer.
func (w LogMiddlewareResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Body will retrieve the stored bytes given on the previous calls
// to the Write method.
func (w LogMiddlewareResponseWriter) Body() []byte {
	return w.body.Bytes()
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareResponseReader
/// ---------------------------------------------------------------------------

// LogMiddlewareResponseReader defines the interface methods of a response
// context reader used to compose the data to be sent to the logger on a
// response event.
type LogMiddlewareResponseReader interface {
	Get(context GinContext) map[string]interface{}
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareBasicResponseReader
/// ---------------------------------------------------------------------------

type LogMiddlewareBasicResponseReader struct{}

// NewLogMiddlewareBasicResponseReader will instantiate a new basic response context reader.
func NewLogMiddlewareBasicResponseReader() *LogMiddlewareBasicResponseReader {
	return &LogMiddlewareBasicResponseReader{}
}

// Get process the context response and return the data to be
// signaled to the logger.
func (r LogMiddlewareBasicResponseReader) Get(context GinContext) map[string]interface{} {
	response := context.(*gin.Context).Writer.(*LogMiddlewareResponseWriter)

	return map[string]interface{}{
		"status":  response.Status(),
		"headers": r.headers(response),
		"body":    string(response.Body()),
		"time":    time.Now().Format("2006-01-02T15:04:05.000-0700"),
	}
}

func (LogMiddlewareBasicResponseReader) headers(response gin.ResponseWriter) map[string][]string {
	headers := map[string][]string{}
	for index, header := range response.Header() {
		headers[index] = header
	}
	return headers
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareJsonResponseReader
/// ---------------------------------------------------------------------------

type LogMiddlewareJsonResponseReader struct {
	reader LogMiddlewareResponseReader
	model  interface{}
}

// NewLogMiddlewareJsonResponseReader will instantiate a new response event
// context reader JSON decorator used to parse the response   body as a JSON
// and add the parsed content into the logging data.
func NewLogMiddlewareJsonResponseReader(reader LogMiddlewareResponseReader, model interface{}) (*LogMiddlewareJsonResponseReader, error) {
	if reader == nil {
		return nil, fmt.Errorf("invalid nil 'reader' argument")
	}

	return &LogMiddlewareJsonResponseReader{
		reader: reader,
		model:  model,
	}, nil
}

// Get process the context response and add the extra bodyJson if the body
// content can be parsed as JSON.
func (r LogMiddlewareJsonResponseReader) Get(context GinContext) map[string]interface{} {
	data := r.reader.Get(context)

	if err := json.Unmarshal([]byte(data["body"].(string)), &r.model); err == nil {
		data["bodyJson"] = r.model
	}

	return data
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareXmlResponseReader
/// ---------------------------------------------------------------------------

type LogMiddlewareXmlResponseReader struct {
	reader LogMiddlewareResponseReader
	model  interface{}
}

// NewLogMiddlewareXmlResponseReader will instantiate a new response event
// context reader XML decorator used to parse the response body as a XML
// and add the parsed content into the logging data.
func NewLogMiddlewareXmlResponseReader(reader LogMiddlewareResponseReader, model interface{}) (*LogMiddlewareXmlResponseReader, error) {
	if reader == nil {
		return nil, fmt.Errorf("invalid nil 'reader' argument")
	}

	return &LogMiddlewareXmlResponseReader{
		reader: reader,
		model:  model,
	}, nil
}

// Get process the context response and add the extra bodyJson if the body
// content can be parsed as XML.
func (r LogMiddlewareXmlResponseReader) Get(context GinContext) map[string]interface{} {
	data := r.reader.Get(context)

	if err := xml.Unmarshal([]byte(data["body"].(string)), &r.model); err == nil {
		data["bodyXml"] = r.model
	}

	return data
}

/// ---------------------------------------------------------------------------
/// LogMiddlewareParams
/// ---------------------------------------------------------------------------

// LogMiddlewareParams defines the storing structure of the parameters
// used to configure the logging middleware.
type LogMiddlewareParams struct {
	RequestReader      LogMiddlewareRequestReader
	ResponseReader     LogMiddlewareResponseReader
	Next               gin.HandlerFunc
	Logger             *Log
	LogChannel         string
	LogLevel           LogLevel
	LogRequestMessage  string
	LogResponseMessage string
}

// NewLogMiddlewareParameters will instantiate a new log middleware parameters
// instance used to configure a log middleware. If environment variables have
// been set for the log environment, the returned parameters structure will
// reflect those values.
func NewLogMiddlewareParams(next gin.HandlerFunc, logger *Log) (*LogMiddlewareParams, error) {
	if next == nil {
		return nil, fmt.Errorf("invalid nil 'next' argument")
	}
	if logger == nil {
		return nil, fmt.Errorf("invalid nil 'logger' argument")
	}

	logChannel := LogMiddlewareChannel
	if env := os.Getenv(EnvLogMiddlewareChannel); env != "" {
		logChannel = env
	}

	logLevel := LogMiddlewareLevel
	if env := os.Getenv(EnvLogMiddlewareLevel); env != "" {
		env = strings.ToLower(env)
		if l, ok := LogLevelMap[env]; !ok {
			panic(fmt.Errorf("unrecognized logger level : %s", env))
		} else {
			logLevel = l
		}
	}

	logRequestMessage := LogMiddlewareRequestMessage
	if env := os.Getenv(EnvLogMiddlewareRequestMessage); env != "" {
		logRequestMessage = env
	}

	logResponseMessage := LogMiddlewareResponseMessage
	if env := os.Getenv(EnvLogMiddlewareResponseMessage); env != "" {
		logResponseMessage = env
	}

	return &LogMiddlewareParams{
		RequestReader:      NewLogMiddlewareBasicRequestReader(),
		ResponseReader:     NewLogMiddlewareBasicResponseReader(),
		Next:               next,
		Logger:             logger,
		LogChannel:         logChannel,
		LogLevel:           logLevel,
		LogRequestMessage:  logRequestMessage,
		LogResponseMessage: logResponseMessage,
	}, nil
}

/// ---------------------------------------------------------------------------
/// LogMiddleware
/// ---------------------------------------------------------------------------

// NewLogMiddleware will instantiate a new middleware that will emit logging
// signals on a request event and on a response event.
func NewLogMiddleware(p *LogMiddlewareParams) (func(ctx *gin.Context), error) {
	if p == nil {
		return nil, fmt.Errorf("invalid nil 'parameters' argument")
	}

	return func(ctx *gin.Context) {
		ctx.Writer, _ = NewLogMiddlewareResponseWriter(ctx.Writer)

		request := p.RequestReader.Get(ctx)
		_ = p.Logger.Signal(p.LogChannel, p.LogLevel, p.LogRequestMessage, map[string]interface{}{"request": request})

		p.Next(ctx)

		response := p.ResponseReader.Get(ctx)
		_ = p.Logger.Signal(p.LogChannel, p.LogLevel, p.LogResponseMessage, map[string]interface{}{"request": request, "response": response})
	}, nil
}
