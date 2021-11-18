package main

import (
	"context"
	"fmt"
	"github.com/winkyi/CloudNative/httpserver/engine"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	g, _ := errgroup.WithContext(context.Background())
	r_app := engine.New()
	r_app.GET("/", engine.Index)
	r_app.GET("/healthz", engine.Healthz)

	app := &http.Server{
		Addr:    ":9999",
		Handler: r_app,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	g.Go(func() error {
		log.Print("Server Started")
		return app.ListenAndServe()
	})

	// pprof
	g.Go(func() error {
		return http.ListenAndServe(":8001", http.DefaultServeMux)
	})

	<-done
	log.Print("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

	err := g.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
