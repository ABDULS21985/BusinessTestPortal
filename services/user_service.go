// services/user_service.go
package services

import (
	"github.com/ABDULS21985/test-portal/models"
	"github.com/ABDULS21985/test-portal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the contract for user-related business logic
type UserService interface {
	RegisterUser(user *models.User) error
	GetUserProfile(id uuid.UUID) (*models.User, error)
	UpdateUserProfile(user *models.User) error
	DeleteUser(id uuid.UUID) error
}

type userService struct {
	repo repositories.UserRepository
}

// NewUserService creates a new instance of userService
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) RegisterUser(user *models.User) error {
	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.repo.CreateUser(user)
}

func (s *userService) GetUserProfile(id uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *userService) UpdateUserProfile(user *models.User) error {
	// Optionally, handle password updates here
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	return s.repo.UpdateUser(user)
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	return s.repo.DeleteUser(id)
}
