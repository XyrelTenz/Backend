package firebase

import (
	"backend/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type App struct {
	App  *firebase.App
	Auth *auth.Client
}

// VerifyIDToken satisfies the TokenVerifier interface from internal/middleware
func (a *App) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return a.Auth.VerifyIDToken(ctx, idToken)
}

func InitFirebase(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	// Reconstruct JSON for the SDK
	sa := map[string]string{
		"type":                        cfg.Firebase.Type,
		"project_id":                  cfg.Firebase.ProjectID,
		"private_key_id":              cfg.Firebase.PrivateKeyID,
		"private_key":                 cfg.Firebase.PrivateKey,
		"client_email":                cfg.Firebase.ClientEmail,
		"client_id":                   cfg.Firebase.ClientID,
		"auth_uri":                    cfg.Firebase.AuthURI,
		"token_uri":                   cfg.Firebase.TokenURI,
		"auth_provider_x509_cert_url": cfg.Firebase.AuthProviderCertURL,
		"client_x509_cert_url":        cfg.Firebase.ClientX509CertURL,
		"universe_domain":             cfg.Firebase.UniverseDomain,
	}

	saJSON, err := json.Marshal(sa)
	if err != nil {
		return nil, fmt.Errorf("error marshaling service account: %v", err)
	}

	opt := option.WithCredentialsJSON(saJSON)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting auth client: %v", err)
	}

	log.Println("Firebase Admin SDK initialized successfully")
	return &App{
		App:  app,
		Auth: authClient,
	}, nil
}
