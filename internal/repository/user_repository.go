package repository

import (
	"context"
	// "duking/internal/logger"
	"duking/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// Create implements UserRepository.
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users(username, email, password, role,
	 full_name, phone,is_verified,avatar_url )
	VALUES($1,$2,$3,$4,$5,$6,$7,$8) 
	RETURNING user_id, created_at`

	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password,
		user.Role, user.FullName, user.Phone, user.IsVerified, user.AvatarURL).
		Scan(&user.UserID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute SQL request create_Hotel: %w", err)
	}
	return nil
}

// Delete implements UserRepository.
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	query := `SELECT * FROM users WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// GetAll implements UserRepository.
func (r *userRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	query := `SELECT user_id, username, email, password, role,
	 full_name, phone,created_at, is_verified,avatar_url FROM users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*models.User

	for rows.Next() {
		var u *models.User
		err := rows.Scan(&u.UserID, &u.Username, &u.Email, &u.Password, &u.Role, &u.FullName,
			&u.Phone, &u.CreatedAt, &u.IsVerified, &u.AvatarURL)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetByEmail implements UserRepository.
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT 
			user_id, username, email, password, role, full_name, phone, 
			created_at, is_verified, avatar_url
		FROM users 
		WHERE email = $1
	`

	row := r.db.QueryRow(ctx, query, email)

	var user models.User
	err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.FullName,
		&user.Phone,
		&user.CreatedAt,
		&user.IsVerified,
		&user.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID implements UserRepository.
func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	query := `
		SELECT 
			user_id, username, email, password, role, full_name, phone, 
			created_at, is_verified, avatar_url
		FROM users 
		WHERE user_id = $1
	`
	row := r.db.QueryRow(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.FullName,
		&user.Phone,
		&user.CreatedAt,
		&user.IsVerified,
		&user.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update implements UserRepository.
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET
			username = $1,
			email = $2,
			password = $3,
			role = $4,
			full_name = $5,
			phone = $6,
			is_verified = $7,
			avatar_url = $8,
			updated_at = NOW()
		WHERE user_id = $9
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.FullName,
		user.Phone,
		user.IsVerified,
		user.AvatarURL,
		user.UserID,
	)
	return err
}
