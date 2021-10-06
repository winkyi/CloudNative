package engine

import (
	"net/http"
)

// Index 主页
func Index(c *Context) {
	c.SetHeaders(c.R.Header)
	c.SetEnvToResponseHeader("VERSION")
	c.HTML(http.StatusOK, "<h1>Hello index</h1>")
	c.StringNotCode("\n%q", c.W.Header())
	c.Log("访问了主页")
}

// Healthz
func Healthz(c *Context) {
	c.JSON(200, H{
		"200": "connect ok",
	})
	c.Log("发起了心跳")
}
