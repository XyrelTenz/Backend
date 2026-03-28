package driver_repo

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/domain"
)

type sqlDriverRepository struct {
	db *sql.DB
}

func NewSQLDriverRepository(db *sql.DB) domain.DriverRepository {
	return &sqlDriverRepository{
		db: db,
	}
}

func (r *sqlDriverRepository) Create(ctx context.Context, driver *domain.Driver) error {
	query := `
		INSERT INTO drivers (user_id, plate_number, vehicle_type, vehicle_color, license_number, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		driver.UserID, driver.PlateNumber, driver.VehicleType,
		driver.VehicleColor, driver.LicenseNumber, domain.DriverStatusInactive,
	).Scan(&driver.ID, &driver.CreatedAt, &driver.UpdatedAt)
}

func (r *sqlDriverRepository) FindByID(ctx context.Context, id string) (*domain.Driver, error) {
	query := `
		SELECT 
			id, user_id, plate_number, vehicle_type, vehicle_color, 
			license_number, status, 
			ST_Y(current_location::geometry), ST_X(current_location::geometry), 
			created_at, updated_at
		FROM drivers WHERE id = $1 OR user_id = $1
	`
	driver := &domain.Driver{}
	var lat, lng sql.NullFloat64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&driver.ID, &driver.UserID, &driver.PlateNumber, &driver.VehicleType, &driver.VehicleColor,
		&driver.LicenseNumber, &driver.Status, &lat, &lng,
		&driver.CreatedAt, &driver.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("driver not found")
	}
	if lat.Valid {
		driver.CurrentLat = lat.Float64
	}
	if lng.Valid {
		driver.CurrentLng = lng.Float64
	}
	return driver, err
}

func (r *sqlDriverRepository) FindByUserID(
	ctx context.Context,
	userID string,
) (*domain.Driver, error) {
	return r.FindByID(ctx, userID)
}

func (r *sqlDriverRepository) UpdateLocation(
	ctx context.Context,
	driverID string,
	lat, lng float64,
) error {
	query := `
		UPDATE drivers 
		SET current_location = ST_SetSRID(ST_Point($1, $2), 4326), updated_at = now()
		WHERE id = $3 OR user_id = $3
	`
	_, err := r.db.ExecContext(ctx, query, lng, lat, driverID)
	return err
}

func (r *sqlDriverRepository) UpdateStatus(
	ctx context.Context,
	driverID string,
	status domain.DriverStatus,
) error {
	query := `UPDATE drivers SET status = $1, updated_at = now() WHERE id = $2 OR user_id = $2`
	_, err := r.db.ExecContext(ctx, query, status, driverID)
	return err
}
