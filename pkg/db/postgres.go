package db

import (
	"context"
	"duking/internal/config"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пула: %w", err)
	}
	log.Println("successful database startup")
	return pool, nil
}
