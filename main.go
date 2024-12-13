package main

import (
	"log"
	"net/http"

	"github.com/ABDULS21985/test-portal/config"
	"github.com/ABDULS21985/test-portal/controllers"
	"github.com/ABDULS21985/test-portal/migrations"
	"github.com/ABDULS21985/test-portal/repositories"
	"github.com/ABDULS21985/test-portal/routes"
	"github.com/ABDULS21985/test-portal/services"

	"github.com/gorilla/mux"
)

func main() {
	// Load configurations
	cfg := config.LoadConfig()

	// Run migrations (optional)
	migrations.RunMigrations(cfg.DB)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(cfg.DB)
	passwordResetRepo := repositories.NewPasswordResetRepository(cfg.DB)

	// Initialize services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, []byte(cfg.JWTSecret))
	passwordResetService := services.NewPasswordResetService(userRepo, passwordResetRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)
	passwordResetController := controllers.NewPasswordResetController(passwordResetService)

	// Setup router
	router := mux.NewRouter()
	routes.SetupRoutes(router, authController, userController, passwordResetController, []byte(cfg.JWTSecret))

	// Start server
	log.Printf("Server running on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
