## 制作证书

```shell script
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=winkyi.com/O=winkyi"
```