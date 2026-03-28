package domain

import (
	"context"
	"time"
)

type ChatMessage struct {
	ID        string    `json:"id"`
	RideID    string    `json:"ride_id"`
	SenderID  string    `json:"sender_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatRepository interface {
	Create(ctx context.Context, msg *ChatMessage) error
	GetByRideID(ctx context.Context, rideID string) ([]*ChatMessage, error)
}
