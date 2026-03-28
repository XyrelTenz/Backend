package passenger_repo

import (
	"context"
	"database/sql"

	"backend/internal/domain"
)

type sqlInteractionRepository struct {
	db *sql.DB
}

func NewSQLInteractionRepository(db *sql.DB) domain.InteractionRepository {
	return &sqlInteractionRepository{
		db: db,
	}
}

func (r *sqlInteractionRepository) AddSavedPlace(ctx context.Context, p *domain.SavedPlace) error {
	query := `
		INSERT INTO saved_places (user_id, label, address, type, location)
		VALUES ($1, $2, $3, $4, ST_SetSRID(ST_Point($5, $6), 4326))
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query, p.UserID, p.Name, p.Address, p.Type, p.Lng, p.Lat).
		Scan(&p.ID, &p.CreatedAt)
}

func (r *sqlInteractionRepository) GetSavedPlaces(
	ctx context.Context,
	userID string,
) ([]*domain.SavedPlace, error) {
	query := `
		SELECT id, user_id, label, address, type, ST_Y(location::geometry), ST_X(location::geometry), created_at
		FROM saved_places WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var places []*domain.SavedPlace
	for rows.Next() {
		p := &domain.SavedPlace{}
		if err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Name,
			&p.Address,
			&p.Type,
			&p.Lat,
			&p.Lng,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		places = append(places, p)
	}
	return places, nil
}

func (r *sqlInteractionRepository) DeleteSavedPlace(ctx context.Context, id string) error {
	query := `DELETE FROM saved_places WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *sqlInteractionRepository) AddRating(ctx context.Context, rtg *domain.Rating) error {
	query := `
		INSERT INTO ratings (ride_id, from_id, to_id, score, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query, rtg.RideID, rtg.FromID, rtg.ToID, rtg.Score, rtg.Comment).
		Scan(&rtg.ID, &rtg.CreatedAt)
}

func (r *sqlInteractionRepository) GetAverageRating(
	ctx context.Context,
	userID string,
) (float64, error) {
	query := `SELECT COALESCE(AVG(score), 0) FROM ratings WHERE to_id = $1`
	var avg float64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&avg)
	return avg, err
}

func (r *sqlInteractionRepository) CreateNotification(
	ctx context.Context,
	n *domain.Notification,
) error {
	query := `
		INSERT INTO notifications (user_id, title, body)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query, n.UserID, n.Title, n.Body).Scan(&n.ID, &n.CreatedAt)
}

func (r *sqlInteractionRepository) GetUserNotifications(
	ctx context.Context,
	userID string,
) ([]*domain.Notification, error) {
	query := `SELECT id, user_id, title, body, is_read, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		n := &domain.Notification{}
		if err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Title,
			&n.Body,
			&n.IsRead,
			&n.CreatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}
