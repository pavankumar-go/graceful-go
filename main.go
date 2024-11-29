package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	service = "*"
)

var (
	port  = flag.Int("port", 3333, "http port to listen on")
	delay = flag.Duration("delay", 10*time.Second,
		"add delay to code execution (5s,10s,20s)")
	shutdownTimeout = flag.Duration("shutdown-timeout", 10*time.Second,
		"shutdown timeout (5s,5m,5h) before connections are cancelled")
	graceful = flag.Int("graceful", 0, "enable graceful shutdown 1 or 0")
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(*delay)
		w.WriteHeader(200)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("%s listening on 0.0.0.0:%d with %v shutdown timeout", service, *port, *shutdownTimeout)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-stop

	log.Printf("%s shutting down ...\n", service)

	var ctx = context.TODO()

	if *graceful != 1 {
		log.Printf("%s gracefully ...\n", service)
		ctxT, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
		ctx = ctxT
		defer cancel()
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Printf("%s complete\n", service)
}
