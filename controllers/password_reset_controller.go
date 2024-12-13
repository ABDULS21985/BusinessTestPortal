// controllers/password_reset_controller.go
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ABDULS21985/test-portal/services"
	"github.com/ABDULS21985/test-portal/utils"
)

type PasswordResetController struct {
	passwordResetService services.PasswordResetService
}

// NewPasswordResetController creates a new instance of PasswordResetController
func NewPasswordResetController(passwordResetService services.PasswordResetService) *PasswordResetController {
	return &PasswordResetController{
		passwordResetService: passwordResetService,
	}
}

// RequestPasswordReset handles password reset requests
func (c *PasswordResetController) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := c.passwordResetService.GetUserByEmail(request.Email)
	if err != nil {
		// For security, do not reveal whether the email exists
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "If that email exists, a reset link has been sent."})
		return
	}

	token, err := c.passwordResetService.CreatePasswordResetToken(user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not create reset token")
		return
	}

	// TODO: Send the token via email to the user
	response := map[string]interface{}{
		"message": "Password reset token generated.",
		"token":   token,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// ResetPassword handles the actual password reset
func (c *PasswordResetController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := c.passwordResetService.ResetPassword(request.Token, request.NewPassword); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Password has been reset successfully."})
}
