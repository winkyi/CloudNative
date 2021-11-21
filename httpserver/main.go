package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"github.com/winkyi/CloudNative/httpserver/engine"
	"gopkg.in/ini.v1"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var configfile string
	flag.Set("v", "4")
	flag.StringVar(&configfile, "configfile", "httpserver/config/app.ini", "http server config.")
	flag.Parse()
	glog.V(2).Info("准备启动httpserver...")
	r_app := engine.New()
	r_app.GET("/", engine.Index)
	r_app.GET("/healthz", engine.Healthz)

	iniConf := engine.IniConfig{FilePath: configfile}
	config, err := iniConf.Load()
	if err != nil {
		panic("can not load config")
	}

	app := &http.Server{
		Addr:    config.(*ini.File).Section("server").Key("port").String(),
		Handler: r_app,
	}
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
}
