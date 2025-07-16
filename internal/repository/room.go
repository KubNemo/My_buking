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

func (r *roomRepository) GetRoomsByHotelID(ctx context.Context, hotelID int64) ([]models.Room, error) {
	query := `
	SELECT 
			r.room_id, r.hotel_id, r.room_number, r.room_type_id, r.price_per_night, r.is_available, r.created_at, r.updated_at,
            rt.room_type_id, rt.name, rt.description, rt.capacity, rt.base_price, rt.amenities
	FROM rooms r
    JOIN room_types rt ON r.room_type_id = rt.room_type_id
    WHERE r.hotel_id = $1
    ORDER BY r.room_number`

	rows, err := r.db.Query(ctx, query, hotelID)
	if err != nil {
		return nil, fmt.Errorf("roomRepository: failed to get rooms by hotel ID: %w", err)
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		var roomType models.RoomType
		err := rows.Scan(
			&room.RoomID, &room.HotelID, &room.RoomNumber, &room.RoomTypeID, &room.PricePerNight, &room.IsAvailable, &room.CreatedAt, &room.UpdatedAt,
			&roomType.RoomTypeID, &roomType.Name, &roomType.Description, &roomType.Capacity, &roomType.BasePrice, &roomType.Amenities,
		)
		if err != nil {
			return nil, fmt.Errorf("roomRepository: failed to scan room row: %w", err)
		}
		room.RoomType = &roomType
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("roomRepository: error during rows iteration: %w", err)
	}
	return rooms, nil
}

func (r *roomRepository) UpdateRoom(ctx context.Context, room *models.Room) error {
	query := `
        UPDATE rooms
        SET
            room_number = $1,
            room_type_id = $2,
            price_per_night = $3,
            is_available = $4,
            updated_at = NOW()
        WHERE room_id = $5
    `
	commandTag, err := r.db.Exec(ctx, query,
		room.RoomNumber,
		room.RoomTypeID,
		room.PricePerNight,
		room.IsAvailable,
		room.RoomID,
	)
	if err != nil {
		return fmt.Errorf("roomRepository: failed to update room: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows // Номер не найден для обновления
	}
	return nil
}

func (r *roomRepository) DeleteRoom(ctx context.Context, id int64) error {
	query := `DELETE FROM rooms WHERE room_id = $1`
	commandTag, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("roomRepository: failed to delete room: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *roomRepository) CheckRoomAvailability(ctx context.Context, roomID int64, checkIn, checkOut time.Time) (bool, error) {
	// Эта логика будет более сложной и скорее всего, будет находиться в BookingRepository
	// Здесь просто проверяем, существует ли номер и доступен ли он по флагу is_available
	query := `SELECT is_available FROM rooms WHERE room_id = $1`
	var isAvailable bool
	err := r.db.QueryRow(ctx, query, roomID).Scan(&isAvailable)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil // Номер не найден
		}
		return false, fmt.Errorf("roomRepository: failed to check room availability: %w", err)
	}
	return isAvailable, nil
}
