### 模块12作业

#### httpserver 服务以 Istio Ingress Gateway 的形式发布出来



#### 创建sidecar

```
kubectl create ns module12
kubectl label ns module12 istio-injection=enabled
```



```
kubectl apply -f httpserver/deploy/httpserver-cm.yaml -n module12 
kubectl apply -f httpserver/deploy/httpserver.yaml -n module12
```





#### 如何实现安全保证



创建tls证书和credential

```shell
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=winkyi Inc./CN=*.winkyi.io' -keyout winkyi.io.key -out winkyi.io.crt
kubectl create -n istio-system secret tls winkyi-credential --key=winkyi.io.key --cert=winkyi.io.crt
```



查询secret创建情况

```shell
winkyi@k8s-dev:~$ kubectl get secret -n istio-system
NAME                                               TYPE                                  DATA   AGE
... ...
winkyi-credential                                  kubernetes.io/tls                     2      2d20h
```



创建specs

```
winkyi@k8s-dev:~$ kubectl apply -f istio-specs.yaml -n module12
virtualservice.networking.istio.io/httpsserver created
gateway.networking.istio.io/httpsserver created
```



查询ingressgateway的IP地址

```shell
winkyi@k8s-dev:~$ kubectl get svc -n istio-system
NAME                   TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                                                                      AGE
istio-egressgateway    ClusterIP      10.107.255.226   <none>        80/TCP,443/TCP                                                               3d21h
istio-ingressgateway   LoadBalancer   10.96.8.157      <pending>     15021:30487/TCP,80:30376/TCP,443:32210/TCP,31400:30079/TCP,15443:31720/TCP   3d21h
istiod                 ClusterIP      10.105.17.111    <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP                                        3d21h
```



验证https访问httpserver服务

```shell
winkyi@k8s-dev:~$ curl --resolve httpsserver.winkyi.io:443:10.96.8.157 https://httpsserver.winkyi.io/healthz  -k
{"200":"connect ok"}
```



#### 七层路由规则

> 实现灰度发布



##### 创建toolbox

```
kubectl apply -f toolbox.yaml -n module12
```



默认流量转发值v1版本，镜像为v1.4, header带"user: winkyi"的发送至之前写的版本镜像为v1.3

```yaml
  hosts:
    - canary
  http:
    - match:
        - headers:
            user:
              exact: winkyi
      route:
        - destination:
            host: canary
            subset: v2
    - route:
      - destination:
          host: canary
          subset: v1
```



##### 创建httpserver-v2版本

```
kubectl apply -f httpserver-v2.yaml -n module12
```



查询labels

```
winkyi@k8s-dev:~$ kubectl get pod -n module12 -l version=v2
NAME                                    READY   STATUS    RESTARTS   AGE
winkyi-httpserver-v2-78757974cb-tllhr   2/2     Running   0          31s
winkyi@k8s-dev:~$ kubectl get pod -n module12 -l version=v1
NAME                                    READY   STATUS    RESTARTS   AGE
winkyi-httpserver-v1-6fc47ddff8-8mwbj   2/2     Running   0          87s
```



##### 更新istio specs

```
winkyi@k8s-dev:~$  kubectl apply -f istio-specs-canary.yaml -n module12
virtualservice.networking.istio.io/canary created
destinationrule.networking.istio.io/canary created
```



##### 进入toolbox验证

```
[root@toolbox-68f79dd5f8-zlgvk /]# curl winkyi-httpserver:9999/hello -H "user: winkyi"
```





##### 在v2查询日志

```
winkyi@k8s-dev:~$ kubectl logs -f winkyi-httpserver-v2-78757974cb-tllhr  -n module12
I1226 03:23:11.946366       1 main.go:24] 准备启动httpserver...
I1226 03:23:11.946847       1 main.go:47] 服务启动完成...
I1226 08:03:41.653046       1 context.go:101] [127.0.0.6]-0  访问了"/hello"
I1226 08:03:43.958358       1 context.go:101] [127.0.0.6]-0  访问了"/hello"
I1226 08:03:45.122925       1 context.go:101] [127.0.0.6]-0  访问了"/hello"
I1226 08:03:46.525312       1 context.go:101] [127.0.0.6]-0  访问了"/hello"
```



####  open tracing 的接入。



##### 安装tracing组件

```shell
winkyi@k8s-dev:~$ kubectl apply -f deploy/tracing/jaeger.yaml
deployment.apps/jaeger created
service/tracing created
service/zipkin created
service/jaeger-collector created
winkyi@k8s-dev:~$ kubectl get pod -n istio-system
NAME                                    READY   STATUS    RESTARTS   AGE
istio-egressgateway-687f4db598-pplf4    1/1     Running   0          6d20h
istio-ingressgateway-78f69bd5db-lnvg4   1/1     Running   0          6d20h
istiod-76d66d9876-gpf4g                 1/1     Running   0          6d20h
jaeger-5d44bc5c5d-h5q9z                 1/1     Running   0          2m36s

# 更改配置文件
winkyi@k8s-dev:~$ kubectl edit configmap istio -n istio-system
set tracing.sampling=100
```

```yaml
data:
  mesh: |-
    accessLogFile: /dev/stdout
    defaultConfig:
      discoveryAddress: istiod.istio-system.svc:15012
      proxyMetadata: {}
      tracing:
        sampling: 100   #新增此处
        zipkin:
          address: zipkin.istio-system:9411
```



##### 创建serviceA服务，用于httpserver调用

```
kubectl apply -f serviceA/deploy/servera.yaml  -n module12
```



通过网关https访问serviceA服务，交易路径

ingress-gateway(https)   -> httpserver  -> serviceA

```shell
curl --resolve httpsserver.winkyi.io:443:10.96.8.157 https://httpsserver.winkyi.io/serverA  -k
```



##### 启动dashboard

```shell
winkyi@k8s-dev:~$ istioctl dashboard jaeger
http://localhost:16686
Failed to open browser; open http://localhost:16686 in your browser.
```



打开页面查询

![image](https://github.com/winkyi/CloudNative/blob/module12/httpserver/docs/images/tracing-1.JPG)



点击查询详情

![image](https://github.com/winkyi/CloudNative/blob/module12/httpserver/docs/images/tracing-2.JPG)
