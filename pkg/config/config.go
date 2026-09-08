package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	DBDsn            string
	JWTSecret        string
	JWTAccessExpiry  string
	JWTRefreshExpiry string
}

func Load() *Config {
	_ = godotenv.Load()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		dbDsn = "postgresql://mydatabase:ecommerce1234@localhost:5435/dbecommerce?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "ECOMMERCE_SUPERSECRET_KEY"
	}

	return &Config{
		AppPort:          port,
		DBDsn:            dbDsn,
		JWTSecret:        jwtSecret,
		JWTAccessExpiry:  "15m",
		JWTRefreshExpiry: "168h",
	}
}
