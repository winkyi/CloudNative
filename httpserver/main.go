package main

import (
	"context"
	"fmt"
	"github.com/winkyi/CloudNative/httpserver/engine"
	"golang.org/x/sync/errgroup"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	g, _ := errgroup.WithContext(context.Background())
	r_app := engine.New()
	r_app.GET("/", engine.Index)
	g.Go(func() error {
		return r_app.Run(":9999")
	})

	r_health := engine.New()
	r_health.GET("/healthz", engine.Healthz)
	g.Go(func() error {
		return r_health.Run(":80")
	})

	// pprof
	g.Go(func() error {
		return http.ListenAndServe(":8001", http.DefaultServeMux)
	})

	err := g.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
