package dto

import "time"

type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

type MessageResponse struct {
	ID        string    `json:"id"`
	RideID    string    `json:"ride_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
