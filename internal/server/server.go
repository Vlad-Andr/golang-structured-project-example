package server

import (
	"best-structure-example/internal/config"
	"context"
	"net/http"
	"time"
)

// Server represents the HTTP server
type Server struct {
	server *http.Server
}

// NewServer creates a new HTTP server
func NewServer(config config.ServerConfig, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:         config.Addr,
			Handler:      handler,
			ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(config.IdleTimeout) * time.Second,
		},
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
