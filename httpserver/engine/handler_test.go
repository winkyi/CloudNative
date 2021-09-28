package engine

import (
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"
)

// TestIndex 测试获取的response header中是否有VERSION
func TestIndex(t *testing.T) {

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*2)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 2))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 2,
		},
	}

	resp, err := client.Get("http://127.0.0.1:9999/4343")

	if err != nil {
		t.Fatal("请求主页失败")
	}

	respHeaders := resp.Header

	defer resp.Body.Close()
	// 获取body
	_, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal("读取body失败")
	}

	if _, ok := respHeaders["Version"]; !ok {
		t.Fatal("response header中VERSION不存在")
	}

}
