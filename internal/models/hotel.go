package models

import "time"

type Hotel struct {
	HotelID     int64     `json:"hotel_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	ImageURL    string    `json:"image_url"` // сюда будет линк из S3
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
