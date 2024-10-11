package config

import (
	"os"
	// this will automatically load your .env file:
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	DATABASE_URL      string
	JWT_SECRET        string
	GoogleLoginConfig oauth2.Config
}

type Handler struct {
	DB     *pgxpool.Pool
	Config *Config
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
	}

	return cfg, nil
}
