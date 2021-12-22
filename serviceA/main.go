package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("starting serverA")

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	srv := http.Server{
		Addr:    ":8088",
		Handler: mux,
	}

	// 优雅终止
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("listen: %s\n", err)
		}
	}()

	glog.Info("serverA is started")
	<-done
	glog.Info("serverA is stoped")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		glog.Fatalf("serverA shutdown failed:%+v", err)
	}
}

func randInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("call index")
	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	io.WriteString(w, "================http request==================")

	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	glog.Infof("respond in %d ms", delay)
}
