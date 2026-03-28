package config

import (
	"encoding/json"
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

	// Support for single JSON string (easier for Render)
	serviceAccountJSON := getEnv("FIREBASE_SERVICE_ACCOUNT_JSON", "")
	if serviceAccountJSON != "" {
		var sa map[string]string
		if err := json.Unmarshal([]byte(serviceAccountJSON), &sa); err == nil {
			cfg.Firebase.Type = sa["type"]
			cfg.Firebase.ProjectID = sa["project_id"]
			cfg.Firebase.PrivateKeyID = sa["private_key_id"]
			cfg.Firebase.PrivateKey = cleanPrivateKey(sa["private_key"])
			cfg.Firebase.ClientEmail = sa["client_email"]
			cfg.Firebase.ClientID = sa["client_id"]
			cfg.Firebase.AuthURI = sa["auth_uri"]
			cfg.Firebase.TokenURI = sa["token_uri"]
			cfg.Firebase.AuthProviderCertURL = sa["auth_provider_x509_cert_url"]
			cfg.Firebase.ClientX509CertURL = sa["client_x509_cert_url"]
			cfg.Firebase.UniverseDomain = sa["universe_domain"]
		}
	} else {
		// If not using JSON, ensure the individual PRIVATE_KEY is also cleaned
		cfg.Firebase.PrivateKey = cleanPrivateKey(cfg.Firebase.PrivateKey)
	}

	checkRequiredFields(cfg)

	return cfg
}

func cleanPrivateKey(key string) string {
	// Remove accidental leading/trailing quotes
	key = strings.Trim(key, "\"")
	// Convert literal \n to actual newlines
	return strings.ReplaceAll(key, "\\n", "\n")
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
