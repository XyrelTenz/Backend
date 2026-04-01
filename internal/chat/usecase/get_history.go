package usecase

import (
	"context"

	"backend/internal/domain"
)

type GetChatHistoryUsecase interface {
	Execute(ctx context.Context, rideID string) ([]*domain.ChatMessage, error)
}

type getChatHistoryUsecase struct {
	repo domain.ChatRepository
}

func NewGetChatHistoryUsecase(repo domain.ChatRepository) GetChatHistoryUsecase {
	return &getChatHistoryUsecase{
		repo: repo,
	}
}

func (u *getChatHistoryUsecase) Execute(
	ctx context.Context,
	rideID string,
) ([]*domain.ChatMessage, error) {
	return u.repo.GetByRideID(ctx, rideID)
}
