package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"backend/internal/config"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
	db         *sql.DB
}

func NewServer(cfg *config.Config, db *sql.DB) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) Run() error {
	router := NewRouter(s.cfg, s.db)
	s.httpServer = &http.Server{
		Addr:    ":" + s.cfg.Server.Port,
		Handler: router,
	}
	log.Printf("Server starting on port %s", s.cfg.Server.Port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
