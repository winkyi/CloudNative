package main

import (
	"github.com/winkyi/CloudNative/httpserver/engine"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	r_app := engine.New()
	r_app.GET("/", engine.Index)
	go r_app.Run(":9999")

	r_health := engine.New()
	r_health.GET("/healthz", engine.Healthz)
	go r_health.Run(":80")

	// pprof
	go http.ListenAndServe(":8001", http.DefaultServeMux)
	select {}
}
