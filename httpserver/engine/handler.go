package engine

import (
	"github.com/winkyi/CloudNative/httpserver/metrics"
	"net/http"
	"strings"
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

func CallServerA(c *Context) {
	c.Log(3, "call %q\n", c.Pattern)
	delay := RandInt(0, 1000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	c.SetHeaders(c.R.Header)
	// 调用serverA服务
	req, err := http.NewRequest("GET", "http://serviceA:8088", nil)
	if err != nil {
		c.Log(3, "%s", err)
	}
	lowerHeader := make(http.Header)
	//打印header
	for k, v := range c.R.Header {
		lowerHeader[strings.ToLower(k)] = v
	}
	c.Log(1, "headers: %q\n", lowerHeader)
	req.Header = lowerHeader
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		c.Log(1, "http get failed with error: ", err)
	}

	c.Log(1, "recevice httpcode %d, respond in %d ms", resp.StatusCode, delay)
}
