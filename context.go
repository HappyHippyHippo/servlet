package servlet

import (
	"io"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
)

// Context interface for the gin-gonic context object.
type Context interface {
	Copy() *gin.Context
	HandlerName() string
	HandlerNames() []string
	Handler() gin.HandlerFunc
	FullPath() string
	Next()
	IsAborted() bool
	Abort()
	AbortWithStatusJSON(code int, jsonObj interface{})
	AbortWithStatus(code int)
	AbortWithError(code int, err error) *gin.Error
	Error(err error) *gin.Error
	Set(key string, value interface{})
	Get(key string) (value interface{}, exists bool)
	MustGet(key string) interface{}
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	Param(key string) string
	Query(key string) string
	DefaultQuery(key, defaultValue string) string
	GetQuery(key string) (string, bool)
	QueryArray(key string) []string
	GetQueryArray(key string) ([]string, bool)
	QueryMap(key string) map[string]string
	GetQueryMap(key string) (map[string]string, bool)
	PostForm(key string) string
	DefaultPostForm(key, defaultValue string) string
	GetPostForm(key string) (string, bool)
	PostFormArray(key string) []string
	GetPostFormArray(key string) ([]string, bool)
	PostFormMap(key string) map[string]string
	GetPostFormMap(key string) (map[string]string, bool)
	FormFile(name string) (*multipart.FileHeader, error)
	MultipartForm() (*multipart.Form, error)
	SaveUploadedFile(file *multipart.FileHeader, dst string) error
	Bind(obj interface{}) error
	BindJSON(obj interface{}) error
	BindXML(obj interface{}) error
	BindQuery(obj interface{}) error
	BindYAML(obj interface{}) error
	BindHeader(obj interface{}) error
	BindUri(obj interface{}) error
	MustBindWith(obj interface{}, b binding.Binding) error
	ShouldBind(obj interface{}) error
	ShouldBindJSON(obj interface{}) error
	ShouldBindXML(obj interface{}) error
	ShouldBindQuery(obj interface{}) error
	ShouldBindYAML(obj interface{}) error
	ShouldBindHeader(obj interface{}) error
	ShouldBindUri(obj interface{}) error
	ShouldBindWith(obj interface{}, b binding.Binding) error
	ShouldBindBodyWith(obj interface{}, bb binding.BindingBody) error
	ClientIP() string
	ContentType() string
	IsWebsocket() bool
	Status(code int)
	Header(key, value string)
	GetHeader(key string) string
	GetRawData() ([]byte, error)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	Cookie(name string) (string, error)
	Render(code int, r render.Render)
	HTML(code int, name string, obj interface{})
	IndentedJSON(code int, obj interface{})
	SecureJSON(code int, obj interface{})
	JSONP(code int, obj interface{})
	JSON(code int, obj interface{})
	AsciiJSON(code int, obj interface{})
	PureJSON(code int, obj interface{})
	XML(code int, obj interface{})
	YAML(code int, obj interface{})
	ProtoBuf(code int, obj interface{})
	String(code int, format string, values ...interface{})
	Redirect(code int, location string)
	Data(code int, contentType string, data []byte)
	DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string)
	File(filepath string)
	FileAttachment(filepath, filename string)
	SSEvent(name string, message interface{})
	Stream(step func(w io.Writer) bool) bool
	Negotiate(code int, config gin.Negotiate)
	NegotiateFormat(offered ...string) string
	SetAccepted(formats ...string)
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}