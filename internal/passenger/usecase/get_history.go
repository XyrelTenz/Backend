package usecase

import (
	"context"

	"backend/internal/domain"
)

type GetPassengerHistoryUsecase interface {
	Execute(ctx context.Context, passengerID string) ([]*domain.Ride, error)
}

type getPassengerHistoryUsecase struct {
	rideRepo domain.RideRepository
}

func NewGetPassengerHistoryUsecase(rideRepo domain.RideRepository) GetPassengerHistoryUsecase {
	return &getPassengerHistoryUsecase{
		rideRepo: rideRepo,
	}
}

func (u *getPassengerHistoryUsecase) Execute(
	ctx context.Context,
	passengerID string,
) ([]*domain.Ride, error) {
	return u.rideRepo.GetPassengerHistory(passengerID)
}
