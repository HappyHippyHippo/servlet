package servlet

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Engine interface for the gin-gonic engine object.
type Engine interface {
	gin.IRoutes
	Delims(left, right string) *gin.Engine
	SecureJsonPrefix(prefix string) *gin.Engine
	LoadHTMLGlob(pattern string)
	LoadHTMLFiles(files ...string)
	SetHTMLTemplate(templ *template.Template)
	SetFuncMap(funcMap template.FuncMap)
	NoRoute(handlers ...gin.HandlerFunc)
	NoMethod(handlers ...gin.HandlerFunc)
	Routes() gin.RoutesInfo
	Run(addr ...string) error
	RunTLS(addr, certFile, keyFile string) error
	RunUnix(file string) error
	RunFd(fd int) error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	HandleContext(c *gin.Context)
}
