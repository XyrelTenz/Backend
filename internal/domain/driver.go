package domain

import (
	"context"
	"time"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DriverStatus string

const (
	DriverStatusActive   DriverStatus = "active"
	DriverStatusInactive DriverStatus = "inactive"
	DriverStatusOnTrip   DriverStatus = "on_trip"
)

type Driver struct {
	ID            string       `json:"id"`
	UserID        string       `json:"user_id"`
	PlateNumber   string       `json:"plate_number"`
	VehicleType   string       `json:"vehicle_type"`
	VehicleColor  string       `json:"vehicle_color"`
	LicenseNumber string       `json:"license_number"`
	Status        DriverStatus `json:"status"`
	CurrentLat    float64      `json:"current_lat"`
	CurrentLng    float64      `json:"current_lng"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type DriverRepository interface {
	Create(ctx context.Context, driver *Driver) error
	FindByID(ctx context.Context, id string) (*Driver, error)
	FindByUserID(ctx context.Context, userID string) (*Driver, error)
	UpdateLocation(ctx context.Context, driverID string, lat, lng float64) error
	UpdateStatus(ctx context.Context, driverID string, status DriverStatus) error
}
