package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockTokenVerifier is a mock implementation of the TokenVerifier interface.
type MockTokenVerifier struct {
	VerifyFunc func(ctx context.Context, idToken string) (*auth.Token, error)
}

func (m *MockTokenVerifier) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (*auth.Token, error) {
	return m.VerifyFunc(ctx, idToken)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Missing Authorization Header", func(t *testing.T) {
		r := gin.New()
		r.Use(AuthMiddleware(&MockTokenVerifier{}))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Authorization header is required")
	})

	t.Run("Invalid Authorization Format", func(t *testing.T) {
		r := gin.New()
		r.Use(AuthMiddleware(&MockTokenVerifier{}))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "InvalidToken")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid authorization format")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		mockVerifier := &MockTokenVerifier{
			VerifyFunc: func(ctx context.Context, idToken string) (*auth.Token, error) {
				return nil, errors.New("invalid token")
			},
		}

		r := gin.New()
		r.Use(AuthMiddleware(mockVerifier))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid or expired token")
	})

	t.Run("Valid Token", func(t *testing.T) {
		mockVerifier := &MockTokenVerifier{
			VerifyFunc: func(ctx context.Context, idToken string) (*auth.Token, error) {
				return &auth.Token{UID: "test-user-id"}, nil
			},
		}

		r := gin.New()
		r.Use(AuthMiddleware(mockVerifier))
		r.GET("/test", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "test-user-id")
	})
}
