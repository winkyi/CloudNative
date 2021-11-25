

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







## 模块8 作业-2



### 1、 集群高可用

此httpserver为无状态应用，使用deployment方式部署，副本数设置为2

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: winkyi-httpserver
spec:
  selector:
    matchLabels:
      app: httpserver
  replicas: 2
... ...
```



部署后pod和deployment控制器

```shell
winkyi@k8s-dev:~$ kubectl get deployment winkyi-httpserver
NAME                READY   UP-TO-DATE   AVAILABLE   AGE
winkyi-httpserver   2/2     2            2           3d13h
winkyi@k8s-dev:~$ kubectl get pod 
NAME                                 READY   STATUS             RESTARTS          AGE
... ...
winkyi-httpserver-69dd4465df-mlrkf   1/1     Running            0                 3d13h
winkyi-httpserver-69dd4465df-skcd5   1/1     Running            0                 3d13h
```





### 2、service方式发布服务

配置service的type为NodePort

```
winkyi@k8s-dev:~$ kubectl get svc winkyi-httpserver
NAME                TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)                         AGE
winkyi-httpserver   NodePort   10.109.72.243   <none>        9999:30988/TCP,8001:30482/TCP   3d12h
```



本地测试

```
winkyi@k8s-dev:~$ curl http://192.168.0.205:30988/
<h1>Hello index</h1>
map["Accept":["*/*"] "Content-Type":["text/html"] "User-Agent":["curl/7.58.0"] "Version":["test_httpserver_env"]]winkyi@k8s-dev:~$ 
winkyi@k8s-dev:~$ curl http://192.168.0.205:30988/healthz
{"200":"connect ok"}
```



### 3、ingress方式发布服务

> 通过证书保证httpserver的通讯安全



#### 制作证书

```shell
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=winkyi.com/O=winkyi"
```



#### 部署secret文件

根据上一步生成的tls.key和tls.crt生成httpserver-secret.yaml文件

```shell
winkyi@k8s-dev:~$ kubectl apply deploy/httpserver-secret.yaml
```



#### 部署ingress文件

```
winkyi@k8s-dev:~$ kubectl apply deploy/ingress.yaml
```



#### 运行结果

```shell
winkyi@k8s-dev:~$ kubectl get svc winkyi-httpserver
NAME                TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)                         AGE
winkyi-httpserver   NodePort   10.109.72.243   <none>        9999:30988/TCP,8001:30482/TCP   3d12h
winkyi@k8s-dev:~/k8s/yaml/ingress$ kubectl get svc -n ingress-nginx
NAME                                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             NodePort    10.104.128.129   <none>        80:32452/TCP,443:32382/TCP   30h
ingress-nginx-controller-admission   ClusterIP   10.109.113.170   <none>        443/TCP                      30h
winkyi@k8s-dev:~/k8s/yaml/ingress$ curl -H "Host: winkyi.com" https://192.168.0.205:32382 -v -k
* Rebuilt URL to: https://192.168.0.205:32382/
*   Trying 192.168.0.205...
* TCP_NODELAY set
* Connected to 192.168.0.205 (192.168.0.205) port 32382 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/certs/ca-certificates.crt
  CApath: /etc/ssl/certs
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS Unknown, Certificate Status (22):
* TLSv1.3 (IN), TLS handshake, Unknown (8):
* TLSv1.3 (IN), TLS Unknown, Certificate Status (22):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS Unknown, Certificate Status (22):
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
* TLSv1.3 (IN), TLS Unknown, Certificate Status (22):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS change cipher, Client hello (1):
* TLSv1.3 (OUT), TLS Unknown, Certificate Status (22):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS_AES_256_GCM_SHA384
* ALPN, server accepted to use h2
* Server certificate:
*  subject: O=Acme Co; CN=Kubernetes Ingress Controller Fake Certificate
*  start date: Nov 22 03:47:06 2021 GMT
*  expire date: Nov 22 03:47:06 2022 GMT
*  issuer: O=Acme Co; CN=Kubernetes Ingress Controller Fake Certificate
*  SSL certificate verify result: unable to get local issuer certificate (20), continuing anyway.
* Using HTTP2, server supports multi-use
* Connection state changed (HTTP/2 confirmed)
* Copying HTTP/2 data in stream buffer to connection buffer after upgrade: len=0
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* Using Stream ID: 1 (easy handle 0x55ba9f55c600)
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
> GET / HTTP/2
> Host: winkyi.com
> User-Agent: curl/7.58.0
> Accept: */*
> 
* TLSv1.3 (IN), TLS Unknown, Certificate Status (22):
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* TLSv1.3 (IN), TLS Unknown, Certificate Status (22):
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* TLSv1.3 (IN), TLS Unknown, Unknown (23):
* Connection state changed (MAX_CONCURRENT_STREAMS updated)!
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* TLSv1.3 (IN), TLS Unknown, Unknown (23):
< HTTP/2 200 
< date: Tue, 23 Nov 2021 09:58:13 GMT
< content-type: text/html
< content-length: 395
< accept: */*
< user-agent: curl/7.58.0
< version: test_httpserver_env
< x-forwarded-for: 192.168.0.205
< x-forwarded-host: winkyi.com
< x-forwarded-port: 443
< x-forwarded-proto: https
< x-forwarded-scheme: https
< x-real-ip: 192.168.0.205
< x-request-id: 57e239505f44d6230d42db5c6215f389
< x-scheme: https
< strict-transport-security: max-age=15724800; includeSubDomains
< 
<h1>Hello index</h1>
* Connection #0 to host 192.168.0.205 left intact
map["Accept":["*/*"] "Content-Type":["text/html"] "User-Agent":["curl/7.58.0"] "Version":["test_httpserver_env"] "X-Forwarded-For":["192.168.0.205"] "X-Forwarded-Host":["winkyi.com"] "X-Forwarded-Port":["443"] "X-Forwarded-Proto":["https"] "X-Forwarded-Scheme":["https"] "X-Real-Ip":["192.168.0.205"] "X-Request-Id":["57e239505f44d6230d42db5c6215f389"] "X-Scheme":["https"]]
```

