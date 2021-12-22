### 模块12作业

#### httpserver 服务以 Istio Ingress Gateway 的形式发布出来





#### 创建sidecar

```
kubectl create ns module12
kubectl label ns module12 istio-injection=enabled
kubectl create -f httpserver.yaml -n securesvc
```



* 如何实现安全保证；

```
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=winkyi Inc./CN=*.winkyi.io' -keyout winkyi.io.key -out winkyi.io.crt
kubectl create -n istio-system secret tls winkyi-credential --key=winkyi.io.key --cert=winkyi.io.crt
```



* 七层路由规则；





* 考虑 open tracing 的接入。



```
```



