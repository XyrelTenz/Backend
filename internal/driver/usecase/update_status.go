package usecase

import (
	"context"

	"backend/internal/domain"
)

type UpdateRideStatusUsecase interface {
	Execute(ctx context.Context, rideID string, status domain.RideStatus) error
}

type updateRideStatusUsecase struct {
	rideRepo   domain.RideRepository
	driverRepo domain.DriverRepository
}

func NewUpdateRideStatusUsecase(rideRepo domain.RideRepository, driverRepo domain.DriverRepository) UpdateRideStatusUsecase {
	return &updateRideStatusUsecase{
		rideRepo:   rideRepo,
		driverRepo: driverRepo,
	}
}

func (u *updateRideStatusUsecase) Execute(ctx context.Context, rideID string, status domain.RideStatus) error {
	ride, err := u.rideRepo.FindByID(rideID)
	if err != nil {
		return err
	}

	if err := u.rideRepo.UpdateStatus(rideID, status); err != nil {
		return err
	}

	if status == domain.RideStatusCompleted && ride.DriverID != nil {
		return u.driverRepo.UpdateStatus(ctx, *ride.DriverID, domain.DriverStatusActive)
	}

	return nil
}
