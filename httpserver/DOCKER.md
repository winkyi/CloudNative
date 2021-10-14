### 打包httpserver镜像



在上节课中写的Makefile中进行镜像打包操作

```shell
winkyi@k8s-dev:~/go/src/winkyi/CloudNative/httpserver$ make release
echo "building httpserver binary"
building httpserver binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deploy/amd64/ .
echo "building httpserver container"
building httpserver container
docker build -t winkyi/httpserver:v1.0 .
Sending build context to Docker daemon  7.228MB
Step 1/6 : FROM ubuntu:20.04
 ---> 597ce1600cf4
Step 2/6 : ENV VERSION=test_httpserver_env
 ---> Running in dbbb5ea87045
Removing intermediate container dbbb5ea87045
 ---> 5f936699f418
Step 3/6 : LABEL app="httpserver" author="winkyi"
 ---> Running in 7bf4bc74c12a
Removing intermediate container 7bf4bc74c12a
 ---> 3c2fe5a0e6b9
Step 4/6 : ADD deploy/amd64/httpserver /httpserver
 ---> 58ad7a53b30f
Step 5/6 : EXPOSE 80
 ---> Running in bd8a0e79093d
Removing intermediate container bd8a0e79093d
 ---> 250a78ecf0ee
Step 6/6 : ENTRYPOINT /httpserver
 ---> Running in 67a7472e4575
Removing intermediate container 67a7472e4575
 ---> 42664c7c0960
Successfully built 42664c7c0960
Successfully tagged winkyi/httpserver:v1.0
```



查询镜像是否打包完成

```shell
winkyi@k8s-dev:~/go/src/winkyi/CloudNative/httpserver$ docker images |grep winkyi
winkyi/httpserver                                                             v1.0       42664c7c0960   26 seconds ago   80MB
winkyi@k8s-dev:~/go/src/winkyi/CloudNative/httpserver$ 
```



登录dockerhub

```shell
winkyi@k8s-dev:~/go/src/winkyi/CloudNative/httpserver$ docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username: winkyi
Password: 
WARNING! Your password will be stored unencrypted in /home/winkyi/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store

Login Succeeded
```



### 镜像推送至dockerhub



使用docker命令上传镜像

```shell
winkyi@k8s-dev:~/go/src/winkyi/CloudNative/httpserver$ docker push winkyi/httpserver:v1.0
The push refers to repository [docker.io/winkyi/httpserver]
46b1144bc71a: Pushed 
da55b45d310b: Mounted from library/ubuntu 
v1.0: digest: sha256:77ae2a6ce12f547957a7de9518db42b5dd05c186ce2e95b7c59cab04d34338a2 size: 740
```



验证dockerhub中是否存在



* 在https://registry.hub.docker.com/中登录可以查询到



* 在命令行docker search中也能查询

  ```shell
  winkyi@k8s-dev:~/go/src/winkyi/CloudNative/httpserver$ docker search winkyi
  NAME                                        DESCRIPTION                     STARS     OFFICIAL   AUTOMATED
  ... ...                                          
  winkyi/httpserver                                                           0                    
  ... ...
  ```

  



### 使用docker命令启动httpserver

```shell
winkyi@k8s-dev:~$ docker run -itd --name httpserver winkyi/httpserver:v1.0
a4bf04458a8322f6464cde50f2ba9d62ca1e67297b6e2062bb4789bcd334415d
winkyi@k8s-dev:~$ docker ps
CONTAINER ID   IMAGE                    COMMAND                  CREATED         STATUS        PORTS     NAMES
a4bf04458a83   winkyi/httpserver:v1.0   "/bin/sh -c /httpser…"   2 seconds ago   Up 1 second   80/tcp    httpserver
... ...
```



### 通过 nsenter 进入容器查看 IP 配置

```shell
winkyi@k8s-dev:~$ docker inspect httpserver | grep -i pid
            "Pid": 13924,
            "PidMode": "",
            "PidsLimit": null,
winkyi@k8s-dev:~$ nsenter -t 13924 -n ip a
nsenter: cannot open /proc/13924/ns/net: Permission denied
winkyi@k8s-dev:~$ sudo nsenter -t 13924 -n ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
16: eth0@if17: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:03 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.3/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
```



### 验证服务

```shell
winkyi@k8s-dev:~$ curl 172.17.0.3/healthz
{"200":"connect ok"}
winkyi@k8s-dev:~$ curl 172.17.0.3:9999/fdfd
<h1>404 page not found</h1>
map["Accept":["*/*"] "Content-Type":["text/html"] "User-Agent":["curl/7.58.0"] "Version":["test_httpserver_env"]]winkyi@k8s-dev:~$ 
winkyi@k8s-dev:~$ curl 172.17.0.3:9999/
<h1>Hello index</h1>
map["Accept":["*/*"] "Content-Type":["text/html"] "User-Agent":["curl/7.58.0"] "Version":["test_httpserver_env"]]winkyi@k8s-dev:~$ 
```



后端日志

```shell
winkyi@k8s-dev:~$ docker logs httpserver
2021/10/14 07:49:24 [172.17.0.1]-404  404 页面不存在
2021/10/14 07:49:35 [172.17.0.1]-200  发起了心跳
2021/10/14 07:49:44 [172.17.0.1]-404  404 页面不存在
2021/10/14 07:49:51 [172.17.0.1]-200  访问了主页
```

