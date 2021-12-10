# 模块10 作业



## 为 HTTPServer 添加 0-2 秒的随机延时

参数孟老师的httpserver添加0-2秒随机延时

```go
func Hello(c *Context) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	c.Log(3, "访问了%q\n",c.Pattern)
	c.SetHeaders(c.R.Header)
	c.JSON(200, H {
		"200": "hello",
	})
    // 0-2秒随机延时
	time.Sleep(time.Millisecond * time.Duration(RandInt(0,2000)))
}
```



## 为 HTTPServer 项目添加延时 Metric



```go
// PrometheusHandler 转换prometheus的handler为handlerFunc
func PrometheusHandler() HandlerFunc {
	h := promhttp.Handler()

	return func(c *Context) {
		h.ServeHTTP(c.W, c.R)
	}
}
```



```go
func main() {
	var configfile string
	flag.Set("v", "4")
	flag.StringVar(&configfile, "configfile", "httpserver/config/app.ini", "http server config.")
	glog.V(2).Info("准备启动httpserver...")
	metrics.Register()
	r_app := engine.New()
	r_app.GET("/", engine.Index)
	r_app.GET("/healthz", engine.Healthz)
	r_app.GET("/hello", engine.Hello)
	// 注册prometheus handler
	r_app.GET("/metrics", engine.PrometheusHandler())
```





## 将 HTTPServer 部署至测试集群，并完成 Prometheus 配置



```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: winkyi-httpserver
spec:
  template:
    metadata:
      labels:
        app: httpserver
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: app
    spec:
...
        ports:
          - containerPort: 9999
            name: app
          - containerPort: 8001
            name: pprof
```





## 从 Promethus 界面中查询延时指标数据







## （可选）创建一个 Grafana Dashboard 展现延时分配情况
