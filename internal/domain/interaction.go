package domain

import (
	"context"
	"time"
)

type SavedPlace struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"` // Label (e.g. My Home, Work)
	Address   string    `json:"address"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Type      string    `json:"type"` // e.g. home, work, other
	CreatedAt time.Time `json:"created_at"`
}

type Rating struct {
	ID        string    `json:"id"`
	RideID    string    `json:"ride_id"`
	FromID    string    `json:"from_id"`
	ToID      string    `json:"to_id"`
	Score     int       `json:"score"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type InteractionRepository interface {
	AddSavedPlace(ctx context.Context, place *SavedPlace) error
	GetSavedPlaces(ctx context.Context, userID string) ([]*SavedPlace, error)
	DeleteSavedPlace(ctx context.Context, id string) error

	AddRating(ctx context.Context, rating *Rating) error
	GetAverageRating(ctx context.Context, userID string) (float64, error)

	CreateNotification(ctx context.Context, n *Notification) error
	GetUserNotifications(ctx context.Context, userID string) ([]*Notification, error)
}
