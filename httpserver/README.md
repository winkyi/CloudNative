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



注册metrics路由

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



在annotations中添加prometheus相关配置

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
        prometheus.io/port: 9999
    spec:
...
        ports:
          - containerPort: 9999
            name: app
          - containerPort: 8001
            name: pprof
```





## 从 Promethus 界面中查询延时指标数据



安装loki-stack

```shell
helm upgrade --install loki . --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=true,prometheus.server.persistentVolume.storageClass=rook-ceph-block,prometheus.server.persistentVolume.size=1Gi,loki.persistence.enabled=true,loki.persistence.storageClassName=rook-ceph-block,loki.persistence.size=1Gi
```



查询prometheus和grafana部署情况

```shell
winkyi@k8s-dev:~$ kubectl get pod
NAME                                            READY   STATUS    RESTARTS       AGE
loki-0                                          1/1     Running   0              3d6h
loki-grafana-6fc47c4f6d-nflsh                   1/1     Running   0              3d6h
loki-kube-state-metrics-6d8fdf5fd8-27mxq        1/1     Running   0              3d6h
loki-prometheus-alertmanager-58df67b8c6-j5klh   2/2     Running   0              3d6h
loki-prometheus-node-exporter-nwpns             1/1     Running   0              3d6h
loki-prometheus-pushgateway-76668f9bf8-mt4nq    1/1     Running   0              3d6h
loki-prometheus-server-79f99f66b9-c2mtp         2/2     Running   0              3d6h
loki-promtail-59pxp                             1/1     Running   0              3d6h
npd-node-problem-detector-l9s6p                 1/1     Running   1 (3d6h ago)   10d
```



将prometheus和grafana的service类型改成NodePort暴露端口

```shell
winkyi@k8s-dev:~$ kubectl get svc
NAME                            TYPE           CLUSTER-IP       EXTERNAL-IP     PORT(S)             AGE
kubernetes                      ClusterIP      10.96.0.1        <none>          443/TCP             29d
loki                            ClusterIP      10.103.77.160    <none>          3100/TCP            4d23h
loki-grafana                    NodePort       10.105.189.106   <none>          80:31036/TCP        4d23h
loki-headless                   ClusterIP      None             <none>          3100/TCP            4d23h
loki-kube-state-metrics         ClusterIP      10.107.18.239    <none>          8080/TCP            4d23h
loki-prometheus-alertmanager    ClusterIP      10.99.67.33      <none>          80/TCP              4d23h
loki-prometheus-node-exporter   ClusterIP      None             <none>          9100/TCP            4d23h
loki-prometheus-pushgateway     ClusterIP      10.99.99.86      <none>          9091/TCP            4d23h
loki-prometheus-server          NodePort       10.101.61.5      <none>          80:30599/TCP        4d23h
```



使用nodeport打开prometheus页面，查询采集的 ```httpserver_latency_seconds_bucket``` 指标

![image](https://github.com/winkyi/CloudNative/blob/module10/httpserver/docs/images/prometheus.JPG)





## 创建一个 Grafana Dashboard 展现延时分配情况



在grafana页面， 采用json方式导入孟老师绘制的panel

```im
"+" -> "import" -> "Import via panel json"
```



配置每分钟访问httpserver

```shell
winkyi@k8s-dev:~$ kubectl get svc
NAME                            TYPE           CLUSTER-IP       EXTERNAL-IP     PORT(S)             AGE
... ...
winkyi-httpserver               ClusterIP      10.109.72.243    <none>          9999/TCP,8001/TCP   18d
winkyi@k8s-dev:~$ cat curl_hello.sh 
#!/bin/bash
curl http://10.109.72.243:9999/hello
winkyi@k8s-dev:~$ crontab -l
* * * * * sh /home/winkyi/curl_hello.sh
```



在grafana页面查询采集情况

![image](https://github.com/winkyi/CloudNative/blob/module10/httpserver/docs/images/grafana.JPG)
