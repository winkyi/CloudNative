package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type H map[string]interface{}

// Context 上下文信息
type Context struct {
	W http.ResponseWriter
	R *http.Request

	Pattern string
	Method  string

	StatusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:       w,
		R:       r,
		Pattern: r.URL.Path,
		Method:  r.Method,
	}
}

// SetHeader 设置header
func (c *Context) SetHeader(key, value string) {
	c.W.Header().Set(key, value)
}

// Status 设置返回码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.W.WriteHeader(code)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.W.Write([]byte(html))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.W)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.W, err.Error(), 500)
	}
}

// Log 设置日志格式
// [127.0.0.1]-200 xxxxx
func (c *Context) Log(format string, v ...interface{}) {
	log.Printf("[%s]-%d  %v", strings.Split(c.R.RemoteAddr, ":")[0], c.StatusCode, fmt.Sprintf(format, v...))
}
