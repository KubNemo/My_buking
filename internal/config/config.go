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
	DB   DBConfig
	Port string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env") // или "./internal/config/.env" — смотри по расположению
	if err != nil {
		log.Println("⚠️  .env файл не найден, загружаю переменные из окружения")
	}

	cfg := &Config{
		Port: os.Getenv("PORT"),
		DB: DBConfig{
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
			Password: os.Getenv("PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
	}

	log.Printf("Загруженные переменные: %+v\n", cfg.DB)

	return cfg
}
