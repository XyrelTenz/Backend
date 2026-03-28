package server

import (
	"backend/internal/config"
	"backend/pkg/firebase"
	"fmt"
	"log"
)

type Server struct {
	cfg         *config.Config
	firebaseApp *firebase.App
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	// Initialize Firebase
	firebaseApp, err := firebase.InitFirebase(s.cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase: %v", err)
	}
	s.firebaseApp = firebaseApp

	// Setup Router
	r := NewRouter(s.firebaseApp)

	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	return r.Run(addr)
}
