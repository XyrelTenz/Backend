package main

import (
	"backend/internal/config"
	"backend/internal/server"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	srv := server.NewServer(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
