package dto

import "backend/internal/domain"

type UpdateStatusRequest struct {
	Status domain.RideStatus `json:"status" binding:"required"`
}

type UpdateLocationRequest struct {
	Lat float64 `json:"lat" binding:"required"`
	Lng float64 `json:"lng" binding:"required"`
}

type NearbyRidesResponse struct {
	Rides []RideDTO `json:"rides"`
}

type RideDTO struct {
	ID             string  `json:"id"`
	PickupAddress  string  `json:"pickup_address"`
	DropoffAddress string  `json:"dropoff_address"`
	DistanceKM     float64 `json:"distance_km"`
	Fare           float64 `json:"fare"`
	Status         string  `json:"status"`
}
