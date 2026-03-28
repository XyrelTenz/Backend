package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	auth_delivery "backend/internal/auth/delivery/http"
	auth_service "backend/internal/auth/service"
	"backend/internal/auth/usecase"
	"backend/internal/domain"
	"backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)
	u.ID = "user_123"
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

type MockDriverRepository struct {
	mock.Mock
}

func (m *MockDriverRepository) Create(ctx context.Context, d *domain.Driver) error {
	return nil
}

func (m *MockDriverRepository) FindByID(ctx context.Context, id string) (*domain.Driver, error) {
	return nil, nil
}

func (m *MockDriverRepository) FindByUserID(
	ctx context.Context,
	id string,
) (*domain.Driver, error) {
	return nil, nil
}

func (m *MockDriverRepository) UpdateLocation(
	ctx context.Context,
	id string,
	lat, lng float64,
) error {
	return nil
}

func (m *MockDriverRepository) UpdateStatus(
	ctx context.Context,
	id string,
	s domain.DriverStatus,
) error {
	return nil
}

func TestAuthIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Login Successful", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		jwtService := auth_service.NewJWTService("secret")
		loginUC := usecase.NewLoginUsecase(mockUserRepo, jwtService)
		signupUC := usecase.NewSignupUsecase(mockUserRepo, nil, jwtService)
		authController := auth_delivery.NewAuthController(signupUC, loginUC)

		r := gin.New()
		r.POST("/auth/login", authController.Login)

		user := &domain.User{ID: "user_123", Email: "test@example.com"}
		user.HashPassword("password123")

		mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)

		loginReq := gin.H{"email": "test@example.com", "password": "password123"}
		body, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.Response
		json.Unmarshal(w.Body.Bytes(), &resp)

		assert.True(t, resp.Success)
		data := resp.Data.(map[string]interface{})
		assert.NotEmpty(t, data["token"])
	})
}
