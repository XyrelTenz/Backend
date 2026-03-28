package dto

type RequestRideRequest struct {
	PickupAddress  string  `json:"pickup_address"  binding:"required"`
	PickupLat      float64 `json:"pickup_lat"      binding:"required"`
	PickupLng      float64 `json:"pickup_lng"      binding:"required"`
	DropoffAddress string  `json:"dropoff_address" binding:"required"`
	DropoffLat     float64 `json:"dropoff_lat"     binding:"required"`
	DropoffLng     float64 `json:"dropoff_lng"     binding:"required"`
	VehicleType    string  `json:"vehicle_type"    binding:"required"`
	PaymentMethod  string  `json:"payment_method"  binding:"required"`
}

type RideResponse struct {
	ID                    string  `json:"id"`
	PickupAddress         string  `json:"pickup_address"`
	DropoffAddress        string  `json:"dropoff_address"`
	DistanceKM            float64 `json:"distance_km"`
	EstimatedFareAmount   float64 `json:"estimated_fare_amount"`
	EstimatedDurationMins int     `json:"estimated_duration_mins"`
	Status                string  `json:"status"`
	VehicleType           string  `json:"vehicle_type"`
}

type SavedPlaceRequest struct {
	Name    string  `json:"name"    binding:"required"`
	Address string  `json:"address" binding:"required"`
	Lat     float64 `json:"lat"     binding:"required"`
	Lng     float64 `json:"lng"     binding:"required"`
	Type    string  `json:"type"    binding:"required"`
}

type SavedPlaceResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Type    string  `json:"type"`
}
