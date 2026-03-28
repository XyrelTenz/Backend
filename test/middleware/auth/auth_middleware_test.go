package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/auth/service"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID string, role string) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*auth_service.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth_service.Claims), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Missing Authorization Header", func(t *testing.T) {
		mockService := new(MockJWTService)
		r := gin.New()
		r.Use(middleware.AuthMiddleware(mockService))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header is required")
	})

	t.Run("Valid Token", func(t *testing.T) {
		mockService := new(MockJWTService)
		claims := &auth_service.Claims{UserID: "user123"}
		mockService.On("ValidateToken", "valid_token").Return(claims, nil)

		r := gin.New()
		r.Use(middleware.AuthMiddleware(mockService))
		r.GET("/test", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "user123")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		mockService := new(MockJWTService)
		mockService.On("ValidateToken", "invalid_token").Return(nil, assert.AnError)

		r := gin.New()
		r.Use(middleware.AuthMiddleware(mockService))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid token")
	})
}
