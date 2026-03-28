package main

import (
	"database/sql"
	"log"

	"backend/internal/config"
	"backend/internal/server"
	"backend/pkg/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()

	cfg := config.LoadConfig()

	db, err := sql.Open("pgx", cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	app := server.NewServer(cfg, db)

	if err := app.Run(); err != nil {
		logger.Log.Fatal("Server failed to start", zap.Error(err))
	}
}
