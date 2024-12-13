// migrations/migrations.go
package migrations

import (
	"log"

	"github.com/ABDULS21985/test-portal/models"

	"gorm.io/gorm"
)

// RunMigrations applies all database schema updates
func RunMigrations(db *gorm.DB) {
	log.Println("Starting database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.PasswordResetToken{},
		// Add more models here as needed
	)

	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	log.Println("Database migrations completed successfully.")
}
