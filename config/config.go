package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ABDULS21985/test-portal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds all configuration for the application
type Config struct {
	Port            string
	JWTSecret       string
	JWTExpiration   time.Duration
	EncryptionKey   string
	DefaultUserPass string
	DB              *gorm.DB
	SMTPHost        string
	SMTPPort        int
	SMTPUser        string
	SMTPPass        string
	LogLevel        string
	UploadDir       string
	SLAInterval     time.Duration
}

// LoadConfig loads environment variables and initializes the database connection
func LoadConfig() *Config {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found. Using system environment variables.")
	}

	// Parse required environment variables
	jwtExpiration, err := time.ParseDuration(getEnv("JWT_EXPIRATION", "72h"))
	if err != nil {
		log.Fatalf("Invalid JWT_EXPIRATION format: %v", err)
	}

	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DATABASE_HOST", "localhost"),
		getEnv("DATABASE_PORT", "5432"),
		getEnv("DATABASE_USER", "postgres"),
		getEnv("DATABASE_PASSWORD", "password"),
		getEnv("DATABASE_NAME", "testportal"),
		getEnv("DATABASE_SSLMODE", "disable"),
	)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automatic migration
	log.Println("Running automatic migrations...")
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}
	log.Println("Automatic migrations completed.")

	// Return the loaded configuration
	return &Config{
		Port:            getEnv("PORT", "8080"),
		JWTSecret:       getEnv("JWT_SECRET", "mysecretkey"),
		JWTExpiration:   jwtExpiration,
		EncryptionKey:   getEnv("ENCRYPTION_KEY", ""),
		DefaultUserPass: getEnv("DEFAULT_USER_PASSWORD", "password"),
		DB:              db,
		SMTPHost:        getEnv("SMTP_HOST", "localhost"),
		SMTPUser:        getEnv("SMTP_USER", ""),
		SMTPPass:        getEnv("SMTP_PASS", ""),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		UploadDir:       getEnv("UPLOAD_DIR", "./uploads"),
	}
}

// getEnv reads an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// parseEnvInt parses an environment variable as an integer, returning a default value if not set
func parseEnvInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	return fmt.Sscanf(value, "%d")
}
