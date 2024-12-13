// controllers/auth_controller.go
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ABDULS21985/test-portal/services"
	"github.com/ABDULS21985/test-portal/utils"
)

type AuthController struct {
	authService services.AuthService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// LoginUser handles user authentication
func (c *AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, token, err := c.authService.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}
