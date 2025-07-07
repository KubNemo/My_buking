package hostes

import (
	"context"
	"duking/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, hotel *models.Hotel) error
	GetOne(ctx context.Context, id uint) (*models.Hotel, error)
	GetAll(ctx context.Context) ([]models.Hotel, error)
	Update(ctx context.Context, id uint, hotel *models.Hotel) (*models.Hotel, error)
	Delete(ctx context.Context, id uint) error
}

// TODO добавить логер
type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, hotel *models.Hotel) error {
	query := `
	INSERT INTO hotels(name,description,location, image_url) 
	VALUES($1,$2,$3,$4)
	RETURNING hotel_id,created_at,updated_at
	`
	err := r.db.QueryRow(
		ctx, query, &hotel.Name, &hotel.Description, &hotel.Location, &hotel.ImageURL).
		Scan(&hotel.HotelID, &hotel.CreatedAt, &hotel.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute SQL request create_Hotel")
	}
	return nil
}

func (r *repository) GetOne(ctx context.Context, id uint) (*models.Hotel, error) {
	query := `SELECT hotel_id, name, description, location, image_url, created_at, updated_at FROM hotels WHERE hotel_id = $1`
	var h models.Hotel
	err := r.db.QueryRow(ctx, query, id).Scan(
		&h.HotelID, &h.Name, &h.Description, &h.Location, &h.ImageURL, &h.CreatedAt, &h.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &h, nil
}
func (r *repository) GetAll(ctx context.Context) ([]models.Hotel, error) {
	query := `SELECT hotel_id, name, description, location, image_url, created_at, updated_at FROM hotels`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hotels []models.Hotel
	for rows.Next() {
		var h models.Hotel
		err := rows.Scan(&h.HotelID, &h.Name, &h.Description, &h.Location, &h.ImageURL, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return nil, err
		}
		hotels = append(hotels, h)
	}
	return hotels, nil
}

func (r *repository) Update(ctx context.Context, id uint, hotel *models.Hotel) (*models.Hotel, error) {
	query := `
		UPDATE hotels
		SET name = $1, description = $2, location = $3, image_url = $4, updated_at = NOW()
		WHERE hotel_id = $5
		RETURNING hotel_id, name, description, location, image_url, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		hotel.Name,
		hotel.Description,
		hotel.Location,
		hotel.ImageURL,
		id,
	).Scan(&hotel.HotelID, &hotel.Name, &hotel.Description, &hotel.Location, &hotel.ImageURL, &hotel.CreatedAt, &hotel.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return hotel, nil
}
func (r *repository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM hotels WHERE hotel_id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
