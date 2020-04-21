package servlet

import (
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
)

// Context interface for the gin-gonic context object.
type Context interface {
	Abort()
	AbortWithError(code int, err error) *gin.Error
	AbortWithStatus(code int)
	AbortWithStatusJSON(code int, jsonObj interface{})
	AsciiJSON(code int, obj interface{})
	Bind(obj interface{}) error
	BindHeader(obj interface{}) error
	BindJSON(obj interface{}) error
	BindQuery(obj interface{}) error
	BindUri(obj interface{}) error
	BindWith(obj interface{}, b binding.Binding) error
	BindXML(obj interface{}) error
	BindYAML(obj interface{}) error
	ClientIP() string
	ContentType() string
	Cookie(name string) (string, error)
	Copy() *gin.Context
	Data(code int, contentType string, data []byte)
	DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string)
	Deadline() (deadline time.Time, ok bool)
	DefaultPostForm(key, defaultValue string) string
	DefaultQuery(key, defaultValue string) string
	Done() <-chan struct{}
	Err() error
	Error(err error) *gin.Error
	File(filepath string)
	FileAttachment(filepath, filename string)
	FileFromFS(filepath string, fs http.FileSystem)
	FormFile(name string) (*multipart.FileHeader, error)
	FullPath() string
	Get(key string) (value interface{}, exists bool)
	GetBool(key string) (b bool)
	GetDuration(key string) (d time.Duration)
	GetFloat64(key string) (f64 float64)
	GetHeader(key string) string
	GetInt(key string) (i int)
	GetInt64(key string) (i64 int64)
	GetPostForm(key string) (string, bool)
	GetPostFormArray(key string) ([]string, bool)
	GetPostFormMap(key string) (map[string]string, bool)
	GetQuery(key string) (string, bool)
	GetQueryArray(key string) ([]string, bool)
	GetQueryMap(key string) (map[string]string, bool)
	GetRawData() ([]byte, error)
	GetString(key string) (s string)
	GetStringMap(key string) (sm map[string]interface{})
	GetStringMapString(key string) (sms map[string]string)
	GetStringMapStringSlice(key string) (smss map[string][]string)
	GetStringSlice(key string) (ss []string)
	GetTime(key string) (t time.Time)
	HTML(code int, name string, obj interface{})
	Handler() gin.HandlerFunc
	HandlerName() string
	HandlerNames() []string
	Header(key, value string)
	IndentedJSON(code int, obj interface{})
	IsAborted() bool
	IsWebsocket() bool
	JSON(code int, obj interface{})
	JSONP(code int, obj interface{})
	MultipartForm() (*multipart.Form, error)
	MustBindWith(obj interface{}, b binding.Binding) error
	MustGet(key string) interface{}
	Negotiate(code int, config gin.Negotiate)
	NegotiateFormat(offered ...string) string
	Next()
	Param(key string) string
	PostForm(key string) string
	PostFormArray(key string) []string
	PostFormMap(key string) map[string]string
	ProtoBuf(code int, obj interface{})
	PureJSON(code int, obj interface{})
	Query(key string) string
	QueryArray(key string) []string
	QueryMap(key string) map[string]string
	Redirect(code int, location string)
	Render(code int, r render.Render)
	SSEvent(name string, message interface{})
	SaveUploadedFile(file *multipart.FileHeader, dst string) error
	SecureJSON(code int, obj interface{})
	Set(key string, value interface{})
	SetAccepted(formats ...string)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	SetSameSite(samesite http.SameSite)
	ShouldBind(obj interface{}) error
	ShouldBindBodyWith(obj interface{}, bb binding.BindingBody) (err error)
	ShouldBindHeader(obj interface{}) error
	ShouldBindJSON(obj interface{}) error
	ShouldBindQuery(obj interface{}) error
	ShouldBindUri(obj interface{}) error
	ShouldBindWith(obj interface{}, b binding.Binding) error
	ShouldBindXML(obj interface{}) error
	ShouldBindYAML(obj interface{}) error
	Status(code int)
	Stream(step func(w io.Writer) bool) bool
	String(code int, format string, values ...interface{})
	Value(key interface{}) interface{}
	XML(code int, obj interface{})
	YAML(code int, obj interface{})
}
