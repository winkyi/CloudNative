package main

import (
	"context"
	"flag"
	"github.com/winkyi/CloudNative/httpserver/engine"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var configfile string
	flag.StringVar(&configfile, "configfile", "httpserver/config/app.ini", "http server config.")
	flag.Parse()
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
		log.Println("服务启动完成...")
		if err := app.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-done
	log.Println("捕获SIGINT或者SIGTERM信号,服务关闭中...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("服务已优雅关闭...")
}
