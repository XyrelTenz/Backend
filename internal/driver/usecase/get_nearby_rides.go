package usecase

import (
	"context"

	"backend/internal/domain"
)

type GetNearbyRidesUsecase interface {
	Execute(ctx context.Context, lat, lng float64) ([]*domain.Ride, error)
}

type getNearbyRidesUsecase struct {
	rideRepo domain.RideRepository
}

func NewGetNearbyRidesUsecase(rideRepo domain.RideRepository) GetNearbyRidesUsecase {
	return &getNearbyRidesUsecase{
		rideRepo: rideRepo,
	}
}

const SearchRadiusMeters = 5000

func (u *getNearbyRidesUsecase) Execute(ctx context.Context, lat, lng float64) ([]*domain.Ride, error) {
	return u.rideRepo.FindNearbyAvailable(lat, lng, SearchRadiusMeters)
}
