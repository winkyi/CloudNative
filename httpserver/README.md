

## 作业运行结果



### 1、 接收客户端 request，并将 request 中带的 header 写入 response header



运行httpserver，在windows的cmd客户端中使用curl进行访问，将response header打印出来

```powershell
C:\Users\winkyi>curl http://127.0.0.1:9999/
<h1>Hello index</h1>
map["Accept":["*/*"] "Content-Type":["text/html"] "User-Agent":["curl/7.55.1"] "Version":["test_httpserver_env"]]
C:\Users\winkyi>curl http://127.0.0.1:9999/
<h1>Hello index</h1>
map["Accept":["*/*"] "Content-Type":["text/html"] "User-Agent":["curl/7.55.1"] "Version":["test_httpserver_env"]]
C:\Users\winkyi>curl http://127.0.0.1:9999/dedede
<h1>404 page not found</h1>
map["Accept":["*/*"] "Content-Type":["text/html"] "User-Agent":["curl/7.55.1"] "Version":["test_httpserver_env"]]
C:\Users\winkyi>
```





### 2、读取当前系统的环境变量中的 VERSION 配置，并写入 response header

* 在作业运行结果1中已将VERSION环境变量打印

* 编写hander_test.go测试模块进行测试，编写httpclient对httpserver进行访问，读取response中的"VERSION"值进行判断

  ```go
  	if _, ok := respHeaders["Version"]; !ok {
  		t.Fatal("response header中VERSION不存在")
  	}
  ```

  

* 单元测试运行结果

  ```
  === RUN   TestIndex
  --- PASS: TestIndex (0.00s)
  PASS
  
  Process finished with exit code 0
  ```

  

### 3、Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出

标准输出结果

```shell
2021/10/06 20:03:24 [127.0.0.1]-200  访问了主页
2021/10/06 20:03:27 [127.0.0.1]-200  访问了主页
2021/10/06 20:03:29 [127.0.0.1]-404  404 页面不存在
2021/10/06 20:06:50 [127.0.0.1]-404  404 页面不存在
```





### 4、当访问 localhost/healthz 时，应返回200

使用curl运行结果

```powershell
C:\Users\winkyi>curl -i localhost/healthz
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 08 Oct 2021 07:57:55 GMT
Content-Length: 21

{"200":"connect ok"}
```



### 待优化

* httpserver中的几个服务，没有处理关闭的事件
