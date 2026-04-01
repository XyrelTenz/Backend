package domain

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	RolePassenger UserRole = "passenger"
	RoleDriver    UserRole = "driver"
	RoleNone      UserRole = "none"
)

type User struct {
	ID           string       `json:"id"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	FullName     string       `json:"full_name"`
	Role         UserRole     `json:"role"`
	IsActive     bool         `json:"is_active"`
	PasswordHash string       `json:"-"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    *time.Time   `json:"deleted_at,omitempty"`
	LastSeenAt   *time.Time   `json:"last_seen_at,omitempty"`
	Profile      *UserProfile `json:"profile,omitempty"`
}

type UserProfile struct {
	UserID              string    `json:"user_id"`
	ProfileImageURL     string    `json:"profile_image_url"`
	Country             string    `json:"country"`
	Address             string    `json:"address"`
	DOB                 string    `json:"dob"` // Use string for Date (YYYY-MM-DD)
	AverageRating       float64   `json:"average_rating"`
	TotalTripsCompleted int       `json:"total_trips_completed"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (u *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}
