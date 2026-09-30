//go:build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgresql://mydatabase:ecommerce1234@localhost:5435/dbecommerce?sslmode=disable&timezone=Asia/Bangkok"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	query := "ALTER TABLE payments ADD COLUMN IF NOT EXISTS qr_code TEXT;"
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}

	fmt.Println("SUCCESS: column qr_code added or verified on payments table!")
}
