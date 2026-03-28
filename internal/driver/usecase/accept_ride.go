package usecase

import (
	"context"
	"errors"

	"backend/internal/domain"
)

type AcceptRideUsecase interface {
	Execute(ctx context.Context, rideID, driverID string) error
}

type acceptRideUsecase struct {
	rideRepo   domain.RideRepository
	driverRepo domain.DriverRepository
}

func NewAcceptRideUsecase(rideRepo domain.RideRepository, driverRepo domain.DriverRepository) AcceptRideUsecase {
	return &acceptRideUsecase{
		rideRepo:   rideRepo,
		driverRepo: driverRepo,
	}
}

func (u *acceptRideUsecase) Execute(ctx context.Context, rideID, driverID string) error {
	driver, err := u.driverRepo.FindByID(ctx, driverID)
	if err != nil {
		return err
	}
	if driver.Status != domain.DriverStatusActive {
		return errors.New("driver must be online to accept rides")
	}

	if err := u.rideRepo.Accept(rideID, driverID); err != nil {
		return err
	}

	return u.driverRepo.UpdateStatus(ctx, driverID, domain.DriverStatusOnTrip)
}
