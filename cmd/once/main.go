package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/wvoliveira/once/configs"
	"github.com/wvoliveira/once/internal/app/once"
)

func main() {
	cfg := configs.New()

	db, err := once.NewStorageDB(cfg)
	if err != nil {
		slog.Error("error to create new storage DB", "error", err.Error())
		os.Exit(1)
	}

	router := Router()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTPServerPort),
		Handler:      router,
		ReadTimeout:  cfg.HTTPServerReadTimeout,
		WriteTimeout: cfg.HTTPServerWriteTimeout,
	}

	closeIdleConnections := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("Error to HTTP server Shutdown: %v", err)
		}
		close(closeIdleConnections)
	}()

	log.Println("HTTP server http://127.0.0.1:8080")

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	log.Println("Initializing shutdown HTTP server...")
	<-closeIdleConnections
	log.Println("HTTP server closed with success!")
}
