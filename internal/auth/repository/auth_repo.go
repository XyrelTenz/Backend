package auth_repo

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type sqlUserRepository struct {
	db *sql.DB
}

func NewSqlUserRepository(db *sql.DB) domain.UserRepository {
	return &sqlUserRepository{
		db: db,
	}
}

func (r *sqlUserRepository) Create(ctx context.Context, user *domain.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Create User
	userID := user.ID
	if userID == "" {
		// Let PostgreSQL generate UUID if not provided
		query := `
			INSERT INTO users (email, phone, full_name, role, password_hash)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at
		`
		err = tx.QueryRowContext(ctx, query, user.Email, user.Phone, user.FullName, user.Role, user.PasswordHash).Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		query := `
			INSERT INTO users (id, email, phone, full_name, role, password_hash)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING created_at, updated_at
		`
		err = tx.QueryRowContext(ctx, query, user.ID, user.Email, user.Phone, user.FullName, user.Role, user.PasswordHash).Scan(
			&user.CreatedAt, &user.UpdatedAt,
		)
	}

	if err != nil {
		return err
	}

	// 2. Create Profile
	profileQuery := `
		INSERT INTO user_profiles (user_id)
		VALUES ($1)
	`
	_, err = tx.ExecContext(ctx, profileQuery, user.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *sqlUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT u.id, u.email, u.phone, u.full_name, u.role, u.password_hash, u.is_active, u.created_at, u.updated_at,
		       p.profile_image_url, p.country, p.address, p.dob, p.average_rating, p.total_trips_completed
		FROM users u
		LEFT JOIN user_profiles p ON u.id = p.user_id
		WHERE u.email = $1 AND u.deleted_at IS NULL
	`
	user := &domain.User{Profile: &domain.UserProfile{}}
	var dob sql.NullString
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Phone, &user.FullName, &user.Role, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		&user.Profile.ProfileImageURL, &user.Profile.Country, &user.Profile.Address, &dob, &user.Profile.AverageRating, &user.Profile.TotalTripsCompleted,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	if dob.Valid {
		user.Profile.DOB = dob.String
	}

	return user, nil
}

func (r *sqlUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT u.id, u.email, u.phone, u.full_name, u.role, u.password_hash, u.is_active, u.created_at, u.updated_at,
		       p.profile_image_url, p.country, p.address, p.dob, p.average_rating, p.total_trips_completed
		FROM users u
		LEFT JOIN user_profiles p ON u.id = p.user_id
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`
	user := &domain.User{Profile: &domain.UserProfile{}}
	var dob sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Phone, &user.FullName, &user.Role, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		&user.Profile.ProfileImageURL, &user.Profile.Country, &user.Profile.Address, &dob, &user.Profile.AverageRating, &user.Profile.TotalTripsCompleted,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	if dob.Valid {
		user.Profile.DOB = dob.String
	}

	return user, nil
}
