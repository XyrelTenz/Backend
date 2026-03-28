package usecase

import (
	"context"
	"time"

	"backend/internal/chat/ws"
	"backend/internal/domain"
)

type SendMessageUsecase interface {
	Execute(ctx context.Context, rideID, senderID, content string) (*domain.ChatMessage, error)
}

type sendMessageUsecase struct {
	repo domain.ChatRepository
	hub  *chat_ws.Hub
}

func NewSendMessageUsecase(repo domain.ChatRepository, hub *chat_ws.Hub) SendMessageUsecase {
	return &sendMessageUsecase{
		repo: repo,
		hub:  hub,
	}
}

func (u *sendMessageUsecase) Execute(ctx context.Context, rideID, senderID, content string) (*domain.ChatMessage, error) {
	msg := &domain.ChatMessage{
		RideID:   rideID,
		SenderID: senderID,
		Message:  content,
	}

	if err := u.repo.Create(ctx, msg); err != nil {
		return nil, err
	}

	u.hub.Broadcast(rideID, map[string]interface{}{
		"type":      "CHAT_MESSAGE",
		"payload":   msg,
		"timestamp": time.Now(),
	})

	return msg, nil
}
