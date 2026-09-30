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

	rows, err := db.Query("SELECT id FROM orders WHERE status = 'pending' AND created_at < NOW() - ($1 * INTERVAL '1 second')", int64(900))
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	var expiredIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			expiredIDs = append(expiredIDs, id)
		}
	}
	fmt.Printf("Expired IDs: %v\n", expiredIDs)

	for _, id := range expiredIDs {
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("tx error: %v", err)
		}

		res, err := tx.Exec(`UPDATE orders SET status = 'cancelled', updated_at = NOW() WHERE id = $1 AND status = 'pending'`, id)
		if err != nil {
			tx.Rollback()
			log.Fatalf("exec update orders error: %v", err)
		}
		aff, _ := res.RowsAffected()
		fmt.Printf("Order %s updated affected=%d\n", id, aff)

		restoreStockQuery := `
			UPDATE products p
			SET stock = p.stock + oi.quantity, updated_at = NOW()
			FROM order_items oi
			WHERE oi.product_id = p.id AND oi.order_id = $1`
		res2, err := tx.Exec(restoreStockQuery, id)
		if err != nil {
			tx.Rollback()
			log.Fatalf("exec restore stock error: %v", err)
		}
		aff2, _ := res2.RowsAffected()
		fmt.Printf("Order %s restore stock affected=%d\n", id, aff2)

		_, _ = tx.Exec(`UPDATE payments SET status = 'cancelled' WHERE order_id = $1 AND status = 'pending'`, id)

		if err := tx.Commit(); err != nil {
			log.Fatalf("commit error: %v", err)
		}
		fmt.Printf("Order %s committed successfully!\n", id)
	}
}
