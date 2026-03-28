package chat_repo

import (
	"context"
	"database/sql"

	"backend/internal/domain"
)

type sqlChatRepository struct {
	db *sql.DB
}

func NewSqlChatRepository(db *sql.DB) domain.ChatRepository {
	return &sqlChatRepository{
		db: db,
	}
}

func (r *sqlChatRepository) Create(ctx context.Context, msg *domain.ChatMessage) error {
	query := `
		INSERT INTO chat_messages (ride_id, sender_id, message)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query, msg.RideID, msg.SenderID, msg.Message).Scan(&msg.ID, &msg.CreatedAt)
}

func (r *sqlChatRepository) GetByRideID(ctx context.Context, rideID string) ([]*domain.ChatMessage, error) {
	query := `
		SELECT id, ride_id, sender_id, message, created_at
		FROM chat_messages WHERE ride_id = $1 ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, rideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.ChatMessage
	for rows.Next() {
		m := &domain.ChatMessage{}
		if err := rows.Scan(&m.ID, &m.RideID, &m.SenderID, &m.Message, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
