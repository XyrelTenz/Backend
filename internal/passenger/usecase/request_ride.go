package usecase

import (
	"context"

	"backend/internal/domain"
)

type RequestRideInput struct {
	PassengerID    string
	PickupAddress  string
	PickupLat      float64
	PickupLng      float64
	DropoffAddress string
	DropoffLat     float64
	DropoffLng     float64
	VehicleType    string
	PaymentMethod  string
}

type RequestRideUsecase interface {
	Execute(ctx context.Context, input RequestRideInput) (*domain.Ride, error)
}

type requestRideUsecase struct {
	rideRepo domain.RideRepository
}

func NewRequestRideUsecase(rideRepo domain.RideRepository) RequestRideUsecase {
	return &requestRideUsecase{
		rideRepo: rideRepo,
	}
}

func (u *requestRideUsecase) Execute(
	ctx context.Context,
	input RequestRideInput,
) (*domain.Ride, error) {
	ride := &domain.Ride{
		PassengerID:    input.PassengerID,
		PickupAddress:  input.PickupAddress,
		PickupLat:      input.PickupLat,
		PickupLng:      input.PickupLng,
		DropoffAddress: input.DropoffAddress,
		DropoffLat:     input.DropoffLat,
		DropoffLng:     input.DropoffLng,
		Status:         domain.RideStatusRequested,
		VehicleType:    input.VehicleType,
		PaymentMethod:  input.PaymentMethod,
	}

	// Calculate fare using domain logic
	ride.CalculateFare()

	if err := u.rideRepo.Create(ride); err != nil {
		return nil, err
	}

	return ride, nil
}
