package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/llauderesv/go-api-gateway/internal/config"
	"github.com/llauderesv/go-api-gateway/internal/gateway"
)

// // Config contains the server configuration.
// type Config struct {
// 	Port   int
// 	Routes []gateway.Route
// }

// Server represents the API Gateway application.
type Server struct {
	httpServer *http.Server
}

// New creates a new Server instance.
func New(cfg config.Config) (*Server, error) {
	router, err := gateway.NewRouter(cfg.Routes)
	if err != nil {
		log.Fatalf("failed to create gateway router: %v", err)
	}

	s := &Server{
		httpServer: &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.Port),
			Handler:           router,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}

	return s, nil
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	log.Printf("Gateway starting on %s", s.httpServer.Addr)

	err := s.httpServer.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// routes registers all HTTP routes.
func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.healthHandler)
	mux.HandleFunc("GET /api/v1/status", s.statusHandler)
	mux.HandleFunc("GET /slow", s.slowHandler)

	return mux
}

type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status: "UP",
		Time:   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

type StatusResponse struct {
	Message string `json:"message"`
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	resp := StatusResponse{
		Message: "API running",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (s *Server) slowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(15 * time.Second)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("finished"))
}
