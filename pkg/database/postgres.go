package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err := db.Exec("SET timezone = 'Asia/Bangkok'"); err != nil {
		log.Printf("Warning: failed to set database timezone: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully (Timezone: Asia/Bangkok)")
	return db, nil
}
