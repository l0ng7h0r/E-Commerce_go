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
	PhajaySecretKey  string
}

func Load() *Config {
	_ = godotenv.Load()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		dbDsn = "postgresql://mydatabase:ecommerce1234@localhost:5432/dbecommerce?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "ecommerce_supersecret_key"
	}

	phajaySecretKey := os.Getenv("PHAJAY_SECRET_KEY")
	if phajaySecretKey == "" {
		phajaySecretKey = "phajay_secret_key"
	}

	return &Config{
		AppPort:          port,
		DBDsn:            dbDsn,
		JWTSecret:        jwtSecret,
		JWTAccessExpiry:  "15m",
		JWTRefreshExpiry: "168h",
		PhajaySecretKey:  phajaySecretKey,
	}
}
