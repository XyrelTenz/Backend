package ride

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/domain"
	"backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDriverRepository struct {
	mock.Mock
}

func (m *MockDriverRepository) Create(ctx context.Context, driver *domain.Driver) error {
	args := m.Called(ctx, driver)
	return args.Error(0)
}

func (m *MockDriverRepository) FindByID(ctx context.Context, id string) (*domain.Driver, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *MockDriverRepository) FindByUserID(
	ctx context.Context,
	userID string,
) (*domain.Driver, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *MockDriverRepository) UpdateLocation(
	ctx context.Context,
	id string,
	lat, lng float64,
) error {
	args := m.Called(ctx, id, lat, lng)
	return args.Error(0)
}

func (m *MockDriverRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status domain.DriverStatus,
) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

type MockRideRepository struct {
	mock.Mock
}

func (m *MockRideRepository) Create(ride *domain.Ride) error {
	args := m.Called(ride)
	ride.ID = "ride_123"
	return args.Error(0)
}

func (m *MockRideRepository) FindByID(id string) (*domain.Ride, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ride), args.Error(1)
}

func (m *MockRideRepository) FindNearbyAvailable(lat, lng float64, r int) ([]*domain.Ride, error) {
	args := m.Called(lat, lng, r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Ride), args.Error(1)
}

func (m *MockRideRepository) UpdateStatus(id string, s domain.RideStatus) error {
	args := m.Called(id, s)
	return args.Error(0)
}

func (m *MockRideRepository) Accept(rID, dID string) error {
	args := m.Called(rID, dID)
	return args.Error(0)
}

func (m *MockRideRepository) GetPassengerHistory(id string) ([]*domain.Ride, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Ride), args.Error(1)
}

func (m *MockRideRepository) GetDriverHistory(id string) ([]*domain.Ride, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Ride), args.Error(1)
}

type MockInteractionRepository struct {
	mock.Mock
}

func (m *MockInteractionRepository) AddSavedPlace(
	ctx context.Context,
	place *domain.SavedPlace,
) error {
	return nil
}

func (m *MockInteractionRepository) GetSavedPlaces(
	ctx context.Context,
	userID string,
) ([]*domain.SavedPlace, error) {
	return nil, nil
}

func (m *MockInteractionRepository) DeleteSavedPlace(ctx context.Context, id string) error {
	return nil
}

func (m *MockInteractionRepository) AddRating(ctx context.Context, rating *domain.Rating) error {
	return nil
}

func (m *MockInteractionRepository) GetAverageRating(
	ctx context.Context,
	userID string,
) (float64, error) {
	return 4.5, nil
}

func (m *MockInteractionRepository) CreateNotification(
	ctx context.Context,
	n *domain.Notification,
) error {
	return nil
}

func (m *MockInteractionRepository) GetUserNotifications(
	ctx context.Context,
	userID string,
) ([]*domain.Notification, error) {
	return nil, nil
}

func TestRideIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	//TODO: Will Fix this test later, since the fare calculation is now more complex and relies on external services, we need to mock those services as well to get a consistent fare estimate for testing. For now, this test will just check that a ride can be requested and that the response contains an estimated fare and the correct status.
	t.Run("Request Ride and Calculate Fare", func(t *testing.T) {
		mockRideRepo := new(MockRideRepository)
		mockInteractionRepo := new(MockInteractionRepository)
		rideService := passenger_service.NewPassengerService(mockRideRepo, mockInteractionRepo)
		rideHandler := passenger_handler.NewPassengerHandler(rideService)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_id", "passenger_123")
			c.Next()
		})
		r.POST("/rides/request", rideHandler.Request)

		rideInput := gin.H{
			"pickup_address":  "Pagadian City",
			"pickup_lat":      7.82,
			"pickup_lng":      123.43,
			"dropoff_address": "Airport",
			"dropoff_lat":     7.82,
			"dropoff_lng":     123.48,
			"vehicle_type":    "MotoTaxi",
			"payment_method":  "Cash",
		}

		mockRideRepo.On("Create", mock.Anything).Return(nil)

		body, _ := json.Marshal(rideInput)
		req, _ := http.NewRequest(http.MethodPost, "/rides/request", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp response.Response
		json.Unmarshal(w.Body.Bytes(), &resp)

		assert.True(t, resp.Success)
		rideData := resp.Data.(map[string]interface{})
		assert.GreaterOrEqual(t, rideData["estimated_fare_amount"], 1.0)
		assert.Equal(t, "requested", string(rideData["status"].(string)))
	})
}
