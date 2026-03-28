package auth_integration_test

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockVerifier for integration test
type MockVerifier struct{}

func (m *MockVerifier) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if idToken == "valid-token" {
		return &auth.Token{UID: "user-123"}, nil
	}
	return nil, assert.AnError
}

func TestGoogleAuthIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Setup router manually for integration test simulation
	r := gin.Default()
	mockVerifier := &MockVerifier{}
	
	// Use similar structure to server.NewRouter
	userHandler := handlers.NewUserHandler(nil) // We don't need firebaseApp here if we mock the middleware
	
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(mockVerifier))
	{
		api.GET("/profile", userHandler.GetProfile)
	}

	t.Run("Authorized Access", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		resp := httptest.NewRecorder()
		
		r.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "user-123")
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		resp := httptest.NewRecorder()
		
		r.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
