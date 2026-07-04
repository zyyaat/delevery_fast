// Package server provides HTTP server setup with graceful shutdown.
package server

import (
        "context"
        "fmt"
        "log/slog"
        "net/http"
        "time"

        "github.com/food-platform/shared/middleware"
)

// Config holds server configuration.
type Config struct {
        Port            int
        ShutdownTimeout time.Duration
        ReadTimeout     time.Duration
        WriteTimeout    time.Duration
        IdleTimeout     time.Duration
        AllowedOrigins  []string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig(port int) Config {
        return Config{
                Port:            port,
                ShutdownTimeout: 30 * time.Second,
                ReadTimeout:     10 * time.Second,
                WriteTimeout:    30 * time.Second,
                IdleTimeout:     120 * time.Second,
                AllowedOrigins:  []string{"*"},
        }
}

// Server wraps http.Server with graceful shutdown.
type Server struct {
        httpServer *http.Server
        cfg        Config
}

// New creates a new Server with the given handler and config.
// It applies standard middleware: RequestID, Logging, Recovery, CORS.
func New(handler http.Handler, cfg Config) *Server {
        handler = middleware.Chain(
                handler,
                middleware.RequestID,
                middleware.Logging,
                middleware.Recovery,
                middleware.CORS(cfg.AllowedOrigins),
        )

        return &Server{
                httpServer: &http.Server{
                        Addr:         fmt.Sprintf(":%d", cfg.Port),
                        Handler:      handler,
                        ReadTimeout:  cfg.ReadTimeout,
                        WriteTimeout: cfg.WriteTimeout,
                        IdleTimeout:  cfg.IdleTimeout,
                },
                cfg: cfg,
        }
}

// Start begins listening for HTTP requests.
// It blocks until Stop() is called or the server fails.
func (s *Server) Start() error {
        slog.Info("http_server_starting", "addr", s.httpServer.Addr)
        if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                return err
        }
        return nil
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
        slog.Info("http_server_shutting_down", "timeout", s.cfg.ShutdownTimeout)

        ctx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
        defer cancel()

        if err := s.httpServer.Shutdown(ctx); err != nil {
                return err
        }

        slog.Info("http_server_stopped")
        return nil
}
