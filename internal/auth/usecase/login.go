package usecase

import (
	"context"
	"errors"

	"backend/internal/auth/service"
	"backend/internal/domain"
)

type LoginUsecase interface {
	Execute(ctx context.Context, email, password string) (string, error)
}

type loginUsecase struct {
	userRepo   domain.UserRepository
	jwtService auth_service.JWTService
}

func NewLoginUsecase(userRepo domain.UserRepository, jwtService auth_service.JWTService) LoginUsecase {
	return &loginUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (u *loginUsecase) Execute(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return "", errors.New("invalid credentials")
	}

	return u.jwtService.GenerateToken(user.ID, string(user.Role))
}
