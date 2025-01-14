package config

import (
	"os"
	// this will automatically load your .env file:
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	DATABASE_URL      string
	JWT_SECRET        string
	GoogleLoginConfig oauth2.Config
	REDIS_DB_URL      string
	REDIS_PASSWORD    string
	REDIS_DB          int
}

type Handler struct {
	DB     *pgxpool.Pool
	Config *Config
	RDB    *redis.Client
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		GoogleLoginConfig: oauth2.Config{
			RedirectURL:  "http://localhost:4000/users/oauth/google/callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint: google.Endpoint,
		},
		REDIS_DB_URL:   os.Getenv("REDIS_URL"),
		REDIS_PASSWORD: os.Getenv("REDIS_PASSWORD"),
		REDIS_DB:       func() int { v, _ := strconv.Atoi(os.Getenv("REDIS_DB")); return v }(),
	}

	return cfg, nil
}
