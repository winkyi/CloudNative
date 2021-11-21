

## 模块8 作业-1



### 1、优雅启动



配置就绪探针 readinessProbe

```yaml
... ...
        readinessProbe:
          httpGet:
            path: /healthz
            port: 9999
          initialDelaySeconds: 10
          periodSeconds: 5
          successThreshold: 2
... ...
```





### 2、优雅终止



参考孟老师 httpserver写法

```go
... ...
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		glog.V(2).Info("服务启动完成...")
		if err := app.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-done
	glog.V(2).Info("捕获SIGINT或者SIGTERM信号,服务关闭中...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := app.Shutdown(ctx); err != nil {
		glog.Fatalf("Server Shutdown Failed:%+v", err)
	}
	glog.V(2).Info("服务已优雅关闭...")
```



终端1 查询日志，观察服务日志

```shell
winkyi@k8s-dev:~/go/src/github.com/winkyi/CloudNative/httpserver/deploy$ kubectl logs -f winkyi-httpserver-69dd4465df-5bcgz 
I1121 11:19:29.083152       1 main.go:22] 准备启动httpserver...
I1121 11:19:29.083468       1 main.go:41] 服务启动完成...
I1121 11:20:20.572018       1 main.go:48] 捕获SIGINT或者SIGTERM信号,服务关闭中...
I1121 11:20:20.573451       1 main.go:57] 服务已优雅关闭...
```





终端2 停止pod

```shell
winkyi@k8s-dev:~$ kubectl delete pod winkyi-httpserver-69dd4465df-5bcgz
pod "winkyi-httpserver-69dd4465df-5bcgz" deleted
winkyi@k8s-dev:~$ 
```







### 3、资源需求和QOS保证



deployment中配置

```yaml
... ...
        resources:
          limits:
            memory: "500Mi"
            cpu: "1"
          requests:
            memory: "200Mi"
            cpu: "200m"
... ...
```



查询QOS

```shell
winkyi@k8s-dev:~$ kubectl describe pod winkyi-httpserver-69dd4465df-l6kc6
... ...
    Restart Count:  0
    Limits:
      cpu:     1
      memory:  500Mi
    Requests:
      cpu:        200m
      memory:     200Mi
... ...
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   Burstable
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age    From               Message
  ----    ------     ----   ----               -------
  Normal  Scheduled  10m    default-scheduler  Successfully assigned default/winkyi-httpserver-69dd4465df-l6kc6 to k8s-dev
  Normal  Pulled     10m    kubelet            Container image "winkyi/httpserver:v1.2" already present on machine
  Normal  Created    10m    kubelet            Created container httpserver
  Normal  Started    9m59s  kubelet            Started container httpserver
... ...
```

>  查询Qos Class为Burstable



### 4、 探活



配置探活探针livenessProbe

```yaml
... ...
        livenessProbe:
          httpGet:
            path: /healthz
            port: 9999
          initialDelaySeconds: 10
          periodSeconds: 5
... ....
```



### 5、 日志等级

默认日志等级为3， 心跳检查日志等级配置为5

```go
// Index 主页
func Index(c *Context) {
	c.SetHeaders(c.R.Header)
	c.SetEnvToResponseHeader("VERSION")
	c.HTML(http.StatusOK, "<h1>Hello index</h1>")
	c.StringNotCode("\n%q", c.W.Header())
	c.Log(3, "访问了主页") // 访问的日志级别3
}

// Healthz
func Healthz(c *Context) {
	c.JSON(200, H{
		"200": "connect ok",
	})
	c.Log(5, "发起了心跳")  // 心跳日志级别5
}
```



配置了/healthz的livenessProbe，为了不让探针日志一直打印，yaml中设置服务启动日志级别为4，可根据需求在yaml文件中调整日志级别。

```yaml
... ...
    spec:
      terminationGracePeriodSeconds: 20  #grace period
      containers:
      - command:
          - httpserver
          - -v=4  # 日志级别
          - -logtostderr
          - -configfile=/config/app.ini
... ...
```



### 6、 配置和代码分离



代码中读取外置INI配置文件

```go
... ...
	iniConf := engine.IniConfig{FilePath: configfile} // 加载外置配置
	config, err := iniConf.Load()
	if err != nil {
		panic("can not load config")
	}

	app := &http.Server{
		Addr:    config.(*ini.File).Section("server").Key("port").String(),  // 读取INI配置文件port配置值
		Handler: r_app,
	}
... ...
```



配置configmap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-httpserver
data:
  app.ini: |
    version = 1.0.0

    [server]
    port = :9999
```



deployment的yaml配置挂载configmap

```yaml
... ...
        volumeMounts:
          - name: httpserverconf
            mountPath: /config/  #挂载在/config文件夹内
            readOnly: true
      volumes:
        -  name: httpserverconf
           configMap:
             name: cm-httpserver
... ...
```



pod中查询configmap挂载情况

```shell
root@winkyi-httpserver-69dd4465df-l6kc6:/# mount | grep config
/dev/vda1 on /config type ext4 (ro,relatime,errors=remount-ro,data=ordered)
root@winkyi-httpserver-69dd4465df-l6kc6:~# cd /config/
root@winkyi-httpserver-69dd4465df-l6kc6:/config# ls
app.ini
root@winkyi-httpserver-69dd4465df-l6kc6:/config# cat app.ini 
version = 1.0.0

[server]
port = :9999
```

