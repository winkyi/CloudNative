package main

import (
	"github.com/winkyi/CloudNative/httpserver/engine"
)

func main() {
	r := engine.New()
	r.GET("/", engine.Index)
	r.GET("/healthz", engine.Healthz)

	r.Run(":9999")
}
