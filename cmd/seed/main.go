package main

import (
	"log"

	"github.com/l0ng7h0r/ecommerce/internal/repository"
	"github.com/l0ng7h0r/ecommerce/pkg/config"
	"github.com/l0ng7h0r/ecommerce/pkg/database"
	"github.com/l0ng7h0r/ecommerce/pkg/security"
)

func main() {
	cfg := config.Load()
	db, err := database.NewPostgres(cfg.DBDsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	adminEmail := "admin@gmail.com"
	adminPass := "admin@123"

	// Check if admin already exists
	existing, _ := userRepo.GetUserByEmail(adminEmail)
	if existing != nil {
		log.Printf("Admin user %s already exists. Ensuring admin role is assigned...", adminEmail)
		_ = userRepo.AssignRole(existing.ID, "admin")
		log.Println("Admin seed completed successfully!")
		return
	}

	hashedPassword, err := security.HashPassword(adminPass)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	user, err := userRepo.CreateUser(adminEmail, hashedPassword)
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	if err := userRepo.AssignRole(user.ID, "admin"); err != nil {
		log.Fatalf("Failed to assign admin role: %v", err)
	}

	log.Printf("Successfully seeded admin user!")
	log.Printf("Email: %s", adminEmail)
	log.Printf("Password: %s", adminPass)
}
