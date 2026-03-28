package passenger_repo

import (
	"database/sql"
	"errors"

	"backend/internal/domain"
)

type sqlRideRepository struct {
	db *sql.DB
}

func NewSQLRideRepository(db *sql.DB) domain.RideRepository {
	return &sqlRideRepository{
		db: db,
	}
}

func (r *sqlRideRepository) Create(ride *domain.Ride) error {
	query := `
		INSERT INTO rides (
			passenger_id, pickup_address, pickup_location, 
			dropoff_address, dropoff_location, distance_km, 
			estimated_duration_mins, estimated_fare_amount, 
			payment_method, status, vehicle_type
		) VALUES (
			$1, $2, ST_SetSRID(ST_Point($3, $4), 4326), 
			$5, ST_SetSRID(ST_Point($6, $7), 4326), $8, 
			$9, $10, $11, $12, $13
		) RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		ride.PassengerID, ride.PickupAddress, ride.PickupLng, ride.PickupLat,
		ride.DropoffAddress, ride.DropoffLng, ride.DropoffLat, ride.DistanceKM,
		ride.EstimatedDurationMins, ride.EstimatedFareAmount,
		ride.PaymentMethod, ride.Status, ride.VehicleType,
	).Scan(&ride.ID, &ride.CreatedAt, &ride.UpdatedAt)
}

func (r *sqlRideRepository) FindByID(id string) (*domain.Ride, error) {
	query := `
		SELECT 
			id, created_at, updated_at, passenger_id, driver_id, 
			pickup_address, ST_Y(pickup_location::geometry) as pickup_lat, ST_X(pickup_location::geometry) as pickup_lng,
			dropoff_address, ST_Y(dropoff_location::geometry) as dropoff_lat, ST_X(dropoff_location::geometry) as dropoff_lng,
			distance_km, estimated_duration_mins, estimated_fare_amount, final_fare_amount,
			payment_method, payment_status, status, vehicle_type,
			cancelled_at, cancellation_reason, completed_at
		FROM rides WHERE id = $1
	`
	ride := &domain.Ride{}
	err := r.db.QueryRow(query, id).Scan(
		&ride.ID, &ride.CreatedAt, &ride.UpdatedAt, &ride.PassengerID, &ride.DriverID,
		&ride.PickupAddress, &ride.PickupLat, &ride.PickupLng,
		&ride.DropoffAddress, &ride.DropoffLat, &ride.DropoffLng,
		&ride.DistanceKM, &ride.EstimatedDurationMins, &ride.EstimatedFareAmount, &ride.FinalFareAmount,
		&ride.PaymentMethod, &ride.PaymentStatus, &ride.Status, &ride.VehicleType,
		&ride.CancelledAt, &ride.CancellationReason, &ride.CompletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("ride not found")
	}
	return ride, err
}

func (r *sqlRideRepository) FindNearbyAvailable(
	lat, lng float64,
	radiusMeter int,
) ([]*domain.Ride, error) {
	query := `
		SELECT 
			id, created_at, pickup_address, ST_Y(pickup_location::geometry), ST_X(pickup_location::geometry),
			dropoff_address, distance_km, estimated_fare_amount, vehicle_type
		FROM rides 
		WHERE status = 'requested' AND driver_id IS NULL
		AND ST_DWithin(pickup_location, ST_SetSRID(ST_Point($1, $2), 4326), $3)
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, lng, lat, radiusMeter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []*domain.Ride
	for rows.Next() {
		ride := &domain.Ride{}
		if err := rows.Scan(
			&ride.ID, &ride.CreatedAt, &ride.PickupAddress, &ride.PickupLat, &ride.PickupLng,
			&ride.DropoffAddress, &ride.DistanceKM, &ride.EstimatedFareAmount, &ride.VehicleType,
		); err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}
	return rides, nil
}

func (r *sqlRideRepository) UpdateStatus(id string, status domain.RideStatus) error {
	var query string
	if status == domain.RideStatusCompleted {
		query = `UPDATE rides SET status = $1, completed_at = now(), payment_status = 'paid' WHERE id = $2`
	} else if status == domain.RideStatusCancelled {
		query = `UPDATE rides SET status = $1, cancelled_at = now() WHERE id = $2`
	} else {
		query = `UPDATE rides SET status = $1, updated_at = now() WHERE id = $2`
	}
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *sqlRideRepository) Accept(rideID, driverID string) error {
	query := `
		UPDATE rides 
		SET status = 'accepted', driver_id = $1, updated_at = now() 
		WHERE id = $2 AND status = 'requested' AND driver_id IS NULL
	`
	res, err := r.db.Exec(query, driverID, rideID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("ride is no longer available")
	}
	return nil
}

func (r *sqlRideRepository) GetPassengerHistory(passengerID string) ([]*domain.Ride, error) {
	query := `SELECT id, created_at, status, estimated_fare_amount FROM rides WHERE passenger_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, passengerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []*domain.Ride
	for rows.Next() {
		ride := &domain.Ride{}
		if err := rows.Scan(
			&ride.ID,
			&ride.CreatedAt,
			&ride.Status,
			&ride.EstimatedFareAmount,
		); err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}
	return rides, nil
}

func (r *sqlRideRepository) GetDriverHistory(driverID string) ([]*domain.Ride, error) {
	query := `SELECT id, created_at, status, estimated_fare_amount FROM rides WHERE driver_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, driverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []*domain.Ride
	for rows.Next() {
		ride := &domain.Ride{}
		if err := rows.Scan(
			&ride.ID,
			&ride.CreatedAt,
			&ride.Status,
			&ride.EstimatedFareAmount,
		); err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}
	return rides, nil
}
