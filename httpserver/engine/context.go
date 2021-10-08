package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

// GetHeader 获取header
func (c *Context) GetHeader() map[string][]string {
	return c.R.Header
}

// SetHeaders 将所有request中的headers配置到response的header中
func (c *Context) SetHeaders(headers map[string][]string) {
	for k, values := range headers {
		for _, v := range values {
			c.SetHeader(k, v)
		}
	}
}

// SetEnvToResponseHeader 获取系统环境变量,并设置到reponse header
func (c *Context) SetEnvToResponseHeader(key string) {
	c.SetHeader(key, os.Getenv(key))
}

// StringNotCode 返回string格式，不带返回code
func (c *Context) StringNotCode(format string, values ...interface{}) {
	c.W.Write([]byte(fmt.Sprintf(format, values...)))
}

// String 返回string格式，带code
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.W.Write([]byte(fmt.Sprintf(format, values...)))
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
	ipaddr := c.R.RemoteAddr
	if strings.Contains(ipaddr, "]") {
		// ipv6格式
		ipaddr = strings.Split(ipaddr, "]")[0][1:]
	} else {
		// ipv4格式
		ipaddr = strings.Split(c.R.RemoteAddr, ":")[0]
	}
	log.Printf("[%s]-%d  %v", ipaddr, c.StatusCode, fmt.Sprintf(format, v...))
}
