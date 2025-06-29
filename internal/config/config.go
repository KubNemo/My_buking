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
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env файл не найден, загружаю переменные из окружения")
	}

	return &Config{
		Port: os.Getenv("PORT"),
		DB: DBConfig{
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
			Password: os.Getenv("PASSWORD"),
		},
	}
}
