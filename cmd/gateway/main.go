package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/llauderesv/go-api-gateway/internal/config"
	"github.com/llauderesv/go-api-gateway/internal/server"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	log.Printf("loaded configuration: %+v", cfg)
	srv, err := server.New(cfg)

	errChan := make(chan error, 1)

	go func() {
		if err := srv.Start(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {
	case err := <-errChan:
		log.Printf("server error: %v", err)
		return
	case sig := <-signalChan:
		log.Printf("Received signal: %v\n", sig)
		log.Println("shutting down server...")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	} else {
		log.Println("server shutdown complete")
	}

	// serverErrors := make(chan error, 1)
	// // Signal handling for graceful shutdown
	// shutdown := make(chan os.Signal, 1)
	// signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// select {
	// case err := <-serverErrors:
	// 	log.Fatalf("Server error: %v", err)

	// case sig := <-shutdown:
	// 	log.Printf("Signal %v received, shutting down gracefully...\n", sig)

	// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 	defer cancel()

	// 	if err := srv.Shutdown(ctx); err != nil {
	// 		log.Printf("Graceful shutdown failed: %v. Forcing exit.", err)
	// 		srv.Close()
	// 	} else {
	// 		log.Println("Server stopped cleanly")
	// 	}
	// }
}
