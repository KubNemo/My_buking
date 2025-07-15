package repository

import (
	"context"
	"duking/internal/models"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, room *models.Room) error
	GetRoomByID(ctx context.Context, id int64) (*models.Room, error)
	GetRoomsByHotelID(ctx context.Context, hotelID int64) ([]models.Room, error)
	UpdateRoom(ctx context.Context, room *models.Room) error
	DeleteRoom(ctx context.Context, id int64) error
	// CheckRoomAvailability проверяет, свободен ли номер на заданный период
	// (логика будет расширена в разделе бронирований)
	CheckRoomAvailability(ctx context.Context, roomID int64, checkIn, checkOut time.Time) (bool, error)
}

type roomRepository struct {
	db *pgxpool.Pool
}

func NewRooomRepository(db *pgxpool.Pool) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) CreateRoom(ctx context.Context, room *models.Room) error {
	query := `
        INSERT INTO rooms(hotel_id, room_number, room_type_id, price_per_night, is_available)
        VALUES($1, $2, $3, $4, $5)
        RETURNING room_id, created_at, updated_at
    `
	err := r.db.QueryRow(ctx, query,
		room.HotelID,
		room.RoomNumber,
		room.RoomTypeID,
		room.PricePerNight,
		room.IsAvailable,
	).Scan(&room.RoomID, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return fmt.Errorf("roomRepository: failed to create room: %w", err)
	}
	return nil
}

func (r *roomRepository) GetRoomByID(ctx context.Context, id int64) (*models.Room, error) {
	query := `
	SELECTr.room_id, r.hotel_id, r.room_number, r.room_type_id, r.price_per_night, r.is_available, r.created_at, r.updated_at,
        rt.room_type_id, rt.name, rt.description, rt.capacity, rt.base_price, rt.amenities
    FROM rooms r 
	JOIN room_types rt ON r.room_type_id = rt.room_type_id
	WHERE r.room_id = $1`

	var room models.Room
	var roomType models.RoomType
	err := r.db.QueryRow(ctx, query, id).Scan(&room.RoomID, &room.HotelID, &room.RoomNumber, &room.RoomTypeID, &room.PricePerNight, &room.IsAvailable, &room.CreatedAt, &room.UpdatedAt,
		&roomType.RoomTypeID, &roomType.Name, &roomType.Description, &roomType.Capacity, &roomType.BasePrice, &roomType.Amenities,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Номер не найден
		}
		return nil, fmt.Errorf("roomRepository: failed to get room by ID: %w", err)
	}
	room.RoomType = &roomType // Присваиваем информацию о типе номера
	return &room, nil
}
