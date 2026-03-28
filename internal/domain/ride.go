package domain

import (
	"math"
	"time"
)

type RideStatus string

const (
	RideStatusRequested RideStatus = "requested"
	RideStatusAccepted  RideStatus = "accepted"
	RideStatusArrived   RideStatus = "arrived"
	RideStatusPickedUp  RideStatus = "picked_up"
	RideStatusCompleted RideStatus = "completed"
	RideStatusCancelled RideStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type Ride struct {
	ID                    string        `json:"id"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
	PassengerID           string        `json:"passenger_id"`
	DriverID              *string       `json:"driver_id,omitempty"`
	PickupAddress         string        `json:"pickup_address"`
	PickupLat             float64       `json:"pickup_lat"`
	PickupLng             float64       `json:"pickup_lng"`
	DropoffAddress        string        `json:"dropoff_address"`
	DropoffLat            float64       `json:"dropoff_lat"`
	DropoffLng            float64       `json:"dropoff_lng"`
	DistanceKM            float64       `json:"distance_km"`
	EstimatedDurationMins int           `json:"estimated_duration_mins"`
	EstimatedFareAmount   float64       `json:"estimated_fare_amount"`
	FinalFareAmount       *float64      `json:"final_fare_amount,omitempty"`
	PaymentMethod         string        `json:"payment_method"`
	PaymentStatus         PaymentStatus `json:"payment_status"`
	Status                RideStatus    `json:"status"`
	VehicleType           string        `json:"vehicle_type"`
	CancelledAt           *time.Time    `json:"cancelled_at,omitempty"`
	CancellationReason    *string       `json:"cancellation_reason,omitempty"`
	CompletedAt           *time.Time    `json:"completed_at,omitempty"`
}

const (
	FarePerKM   = 20.0
	MinimumFare = 40.0
)

func (r *Ride) CalculateFare() {
	distance := CalculateDistance(r.PickupLat, r.PickupLng, r.DropoffLat, r.DropoffLng)
	r.DistanceKM = distance
	r.EstimatedFareAmount = math.Max(MinimumFare, distance*FarePerKM)
	r.EstimatedDurationMins = int(distance * 3) // Assuming 3 mins per KM on average
}

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in KM
	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

type RideRepository interface {
	Create(ride *Ride) error
	FindByID(id string) (*Ride, error)
	FindNearbyAvailable(lat, lng float64, radiusMeter int) ([]*Ride, error)
	UpdateStatus(id string, status RideStatus) error
	Accept(rideID, driverID string) error
	GetPassengerHistory(passengerID string) ([]*Ride, error)
	GetDriverHistory(driverID string) ([]*Ride, error)
}
