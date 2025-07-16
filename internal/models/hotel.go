package models

import "time"

type Hotel struct {
	HotelID     int64     `json:"hotel_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`  // ул. Ленина, д. 25
	City        string    `json:"city"`     // Москва
	State       string    `json:"state"`    // Московская область
	Country     string    `json:"country"`  // Россия
	ZipCode     string    `json:"zip_code"` // 101000
	Latitude    float64   `json:"latitude"` // для карты и геопоиска
	Longitude   float64   `json:"longitude"`
	Stars       int       `json:"stars"`     // от 1 до 5
	Amenities   []string  `json:"amenities"` // ["Wi-Fi", "Парковка", "Бассейн"]
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
