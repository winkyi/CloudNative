package engine

import (
	"github.com/winkyi/CloudNative/httpserver/metrics"
	"net/http"
	"time"
)

// Index 主页
func Index(c *Context) {
	c.SetHeaders(c.R.Header)
	c.SetEnvToResponseHeader("VERSION")
	c.HTML(http.StatusOK, "<h1>Hello index</h1>")
	c.StringNotCode("\n%q", c.W.Header())
	c.Log(3, "访问了主页")
}

// Healthz
func Healthz(c *Context) {
	c.JSON(200, H{
		"200": "connect ok",
	})
	c.Log(5, "发起了心跳")
}

func Hello(c *Context) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	c.Log(3, "访问了%q\n", c.Pattern)
	c.SetHeaders(c.R.Header)
	c.JSON(200, H{
		"200": "hello",
	})
	// 0-2秒随机延时
	time.Sleep(time.Millisecond * time.Duration(RandInt(0, 2000)))
}
