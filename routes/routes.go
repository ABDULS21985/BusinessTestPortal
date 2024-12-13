package routes

import (
	"net/http"

	"github.com/ABDULS21985/test-portal/controllers"
	"github.com/ABDULS21985/test-portal/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes initializes all API routes and associates them with their respective controllers
func SetupRoutes(router *mux.Router, authController *controllers.AuthController, userController *controllers.UserController, passwordResetController *controllers.PasswordResetController, authMiddleware *middleware.AuthMiddleware) {
	api := router.PathPrefix("/api").Subrouter()

	// Auth Routes (Public)
	api.HandleFunc("/auth/login", authController.LoginUser).Methods("POST")

	// User Routes
	// Registration is handled by User Controller
	api.HandleFunc("/users/register", userController.RegisterUser).Methods("POST")

	// Password Reset Routes
	api.HandleFunc("/password-reset/request", passwordResetController.RequestPasswordReset).Methods("POST")
	api.HandleFunc("/password-reset/reset", passwordResetController.ResetPassword).Methods("POST")

	// Protected User Routes
	protected := api.PathPrefix("/users").Subrouter()
	protected.Use(authMiddleware.RequireAuth)
	protected.HandleFunc("/{id}", userController.GetUserProfile).Methods("GET")
	protected.HandleFunc("/{id}", userController.UpdateUserProfile).Methods("PUT")
	protected.HandleFunc("/{id}", userController.DeleteUser).Methods("DELETE")

	// Additional protected routes
	admin := router.PathPrefix("/protected").Subrouter()
	admin.Use(authMiddleware.RequireAuth)
	admin.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome, Admin!"))
	}).Methods("GET")
}
