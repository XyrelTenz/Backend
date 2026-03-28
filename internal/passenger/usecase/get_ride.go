package usecase

import (
	"context"

	"backend/internal/domain"
)

type GetRideUsecase interface {
	Execute(ctx context.Context, id string) (*domain.Ride, error)
}

type getRideUsecase struct {
	rideRepo domain.RideRepository
}

func NewGetRideUsecase(rideRepo domain.RideRepository) GetRideUsecase {
	return &getRideUsecase{
		rideRepo: rideRepo,
	}
}

func (u *getRideUsecase) Execute(ctx context.Context, id string) (*domain.Ride, error) {
	return u.rideRepo.FindByID(id)
}
