package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Auth struct {
		JWTSecret string
	}
	Database struct {
		URL string
	}
	Server struct {
		Port string
	}
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	cfg := &Config{}
	cfg.Auth.JWTSecret = getEnv("JWT_SECRET", "your-secret-key")
	cfg.Database.URL = getEnv("DATABASE_URL", "")
	cfg.Server.Port = getEnv("PORT", "8080")

	checkRequiredFields(cfg)

	return cfg
}

func (c *Config) GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func checkRequiredFields(cfg *Config) {
	required := map[string]string{
		"JWT_SECRET":   cfg.Auth.JWTSecret,
		"DATABASE_URL": cfg.Database.URL,
	}

	for key, val := range required {
		if val == "" {
			log.Printf("WARNING: Environment variable %s is missing or empty!", key)
		}
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
