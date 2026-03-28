package usecase

import (
	"context"

	"backend/internal/domain"
)

type AddSavedPlaceUsecase interface {
	Execute(ctx context.Context, place *domain.SavedPlace) error
}

type addSavedPlaceUsecase struct {
	interactionRepo domain.InteractionRepository
}

func NewAddSavedPlaceUsecase(interactionRepo domain.InteractionRepository) AddSavedPlaceUsecase {
	return &addSavedPlaceUsecase{
		interactionRepo: interactionRepo,
	}
}

func (u *addSavedPlaceUsecase) Execute(ctx context.Context, place *domain.SavedPlace) error {
	return u.interactionRepo.AddSavedPlace(ctx, place)
}
