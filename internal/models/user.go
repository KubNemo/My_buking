package models

import "time"

type User struct {
	UserID     uint      `json:"user_id" db:"user_id"`
	Username   string    `json:"username" db:"username"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"-" db:"password"` // скрыт из JSON
	Role       string    `json:"role" db:"role"`  // admin, user
	FullName   string    `json:"full_name" db:"full_name"`
	Phone      string    `json:"phone,omitempty" db:"phone"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	IsVerified bool      `json:"is_verified" db:"is_verified"` // подтвердил email
	AvatarURL  string    `json:"avatar_url,omitempty" db:"avatar_url"`
}
