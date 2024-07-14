package config

import (
	"os"
	// this will automatically load your .env file:
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DATABASE_URL string
	JWT_SECRET   string
}

type Handler struct {
	DB     *pgxpool.Pool
	Config *Config
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
	}

	return cfg, nil
}
