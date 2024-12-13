// services/password_reset_service.go
package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ABDULS21985/test-portal/models"
	"github.com/ABDULS21985/test-portal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// PasswordResetService defines the contract for password reset operations
type PasswordResetService interface {
	CreatePasswordResetToken(userID uuid.UUID) (string, error)
	ValidatePasswordResetToken(token string) (*models.PasswordResetToken, error)
	ResetPassword(token, newPassword string) error
	GetUserByEmail(email string) (*models.User, error)
}

type passwordResetService struct {
	userRepo          repositories.UserRepository
	passwordResetRepo repositories.PasswordResetRepository
}

// NewPasswordResetService creates a new instance of passwordResetService
func NewPasswordResetService(userRepo repositories.UserRepository, passwordResetRepo repositories.PasswordResetRepository) PasswordResetService {
	return &passwordResetService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
	}
}

func (s *passwordResetService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *passwordResetService) CreatePasswordResetToken(userID uuid.UUID) (string, error) {
	// Generate a secure random token
	tokenBytes := make([]byte, 16)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	prt := &models.PasswordResetToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Token valid for 1 hour
	}

	if err := s.passwordResetRepo.CreateToken(prt); err != nil {
		return "", err
	}

	return token, nil
}

func (s *passwordResetService) ValidatePasswordResetToken(token string) (*models.PasswordResetToken, error) {
	prt, err := s.passwordResetRepo.GetToken(token)
	if err != nil {
		return nil, err
	}

	if prt.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return prt, nil
}

func (s *passwordResetService) ResetPassword(token, newPassword string) error {
	prt, err := s.ValidatePasswordResetToken(token)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByID(prt.UserID)
	if err != nil {
		return err
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Update user password
	if err := s.userRepo.UpdateUser(user); err != nil {
		return err
	}

	// Delete the used token
	if err := s.passwordResetRepo.DeleteToken(prt.ID); err != nil {
		return err
	}

	return nil
}
