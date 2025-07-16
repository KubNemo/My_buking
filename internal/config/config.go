package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Name     string
	Port     string
	Password string
	Host     string
	User     string
	SSLMode  string
}

type Config struct {
	DB        DBConfig
	Port      string
	SekretKey string
	IsProd    bool
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env файл не найден, загружаю переменные из окружения")
	}
	isProd := os.Getenv("APP_ENV") == "production" || os.Getenv("APP_ENV") == "prod"
	return &Config{
		Port:      os.Getenv("PORT"),
		SekretKey: os.Getenv("JWT_SECRET_KEY"),
		IsProd:    isProd,
		DB: DBConfig{
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
			Password: os.Getenv("PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
	}
}
