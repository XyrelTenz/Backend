package usecase

import (
	"context"

	auth_service "backend/internal/auth/service"
	"backend/internal/domain"
)

type SignupInput struct {
	Email         string
	Password      string
	FullName      string
	Phone         string
	Role          domain.UserRole
	PlateNumber   string
	VehicleType   string
	VehicleColor  string
	LicenseNumber string
}

type SignupUsecase interface {
	Execute(ctx context.Context, input SignupInput) (string, error)
}

type signupUsecase struct {
	userRepo   domain.UserRepository
	driverRepo domain.DriverRepository
	jwtService auth_service.JWTService
}

func NewSignupUsecase(
	userRepo domain.UserRepository,
	driverRepo domain.DriverRepository,
	jwtService auth_service.JWTService,
) SignupUsecase {
	return &signupUsecase{
		userRepo:   userRepo,
		driverRepo: driverRepo,
		jwtService: jwtService,
	}
}

func (u *signupUsecase) Execute(ctx context.Context, input SignupInput) (string, error) {
	user := &domain.User{
		Email:    input.Email,
		FullName: input.FullName,
		Phone:    input.Phone,
		Role:     input.Role,
	}

	if err := user.HashPassword(input.Password); err != nil {
		return "", err
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return "", err
	}

	if input.Role == domain.RoleDriver {
		driver := &domain.Driver{
			UserID:        user.ID,
			PlateNumber:   input.PlateNumber,
			VehicleType:   input.VehicleType,
			VehicleColor:  input.VehicleColor,
			LicenseNumber: input.LicenseNumber,
			Status:        domain.DriverStatusInactive,
		}
		if err := u.driverRepo.Create(ctx, driver); err != nil {
			return "", err
		}
	}

	return u.jwtService.GenerateToken(user.ID, string(user.Role))
}
