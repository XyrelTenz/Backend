package usecase

import (
	"context"

	"backend/internal/domain"
)

type UpdateLocationUsecase interface {
	Execute(ctx context.Context, driverID string, lat, lng float64) error
}

type updateLocationUsecase struct {
	driverRepo domain.DriverRepository
}

func NewUpdateLocationUsecase(driverRepo domain.DriverRepository) UpdateLocationUsecase {
	return &updateLocationUsecase{
		driverRepo: driverRepo,
	}
}

func (u *updateLocationUsecase) Execute(ctx context.Context, driverID string, lat, lng float64) error {
	return u.driverRepo.UpdateLocation(ctx, driverID, lat, lng)
}
