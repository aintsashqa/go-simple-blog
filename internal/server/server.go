package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aintsashqa/go-simple-blog/internal/config"
)

type Server struct {
	http *http.Server
}

func NewServer(cfg config.AppConfig, handler http.Handler) *Server {
	return &Server{
		http: &http.Server{
			Addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:        handler,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: cfg.MaxHeaderMBytes << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.http.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
