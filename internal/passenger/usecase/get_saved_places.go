package usecase

import (
	"context"

	"backend/internal/domain"
)

type GetSavedPlacesUsecase interface {
	Execute(ctx context.Context, userID string) ([]*domain.SavedPlace, error)
}

type getSavedPlacesUsecase struct {
	interactionRepo domain.InteractionRepository
}

func NewGetSavedPlacesUsecase(interactionRepo domain.InteractionRepository) GetSavedPlacesUsecase {
	return &getSavedPlacesUsecase{
		interactionRepo: interactionRepo,
	}
}

func (u *getSavedPlacesUsecase) Execute(ctx context.Context, userID string) ([]*domain.SavedPlace, error) {
	return u.interactionRepo.GetSavedPlaces(ctx, userID)
}
