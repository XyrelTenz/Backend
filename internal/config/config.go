package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Firebase struct {
		Type                   string
		ProjectID              string
		PrivateKeyID           string
		PrivateKey             string
		ClientEmail            string
		ClientID               string
		AuthURI                string
		TokenURI               string
		AuthProviderCertURL    string
		ClientX509CertURL      string
		UniverseDomain         string
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
	cfg.Firebase.Type = getEnv("FIREBASE_TYPE", "")
	cfg.Firebase.ProjectID = getEnv("FIREBASE_PROJECT_ID", "")
	cfg.Firebase.PrivateKeyID = getEnv("FIREBASE_PRIVATE_KEY_ID", "")
	
	// Handle private key with escaped newlines
	privateKey := getEnv("FIREBASE_PRIVATE_KEY", "")
	cfg.Firebase.PrivateKey = strings.ReplaceAll(privateKey, "\\n", "\n")
	
	cfg.Firebase.ClientEmail = getEnv("FIREBASE_CLIENT_EMAIL", "")
	cfg.Firebase.ClientID = getEnv("FIREBASE_CLIENT_ID", "")
	cfg.Firebase.AuthURI = getEnv("FIREBASE_AUTH_URI", "")
	cfg.Firebase.TokenURI = getEnv("FIREBASE_TOKEN_URI", "")
	cfg.Firebase.AuthProviderCertURL = getEnv("FIREBASE_AUTH_CERT_URL", "")
	cfg.Firebase.ClientX509CertURL = getEnv("FIREBASE_CLIENT_CERT_URL", "")
	cfg.Firebase.UniverseDomain = getEnv("FIREBASE_UNIVERSE_DOMAIN", "")

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
		"FIREBASE_TYPE":           cfg.Firebase.Type,
		"FIREBASE_PROJECT_ID":     cfg.Firebase.ProjectID,
		"FIREBASE_PRIVATE_KEY":    cfg.Firebase.PrivateKey,
		"FIREBASE_CLIENT_EMAIL":   cfg.Firebase.ClientEmail,
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
