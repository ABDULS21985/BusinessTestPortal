// repositories/password_reset_repository.go
package repositories

import (
	"github.com/ABDULS21985/test-portal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PasswordResetRepository defines the contract for password reset token operations
type PasswordResetRepository interface {
	CreateToken(token *models.PasswordResetToken) error
	GetToken(token string) (*models.PasswordResetToken, error)
	DeleteToken(id uuid.UUID) error
}

// passwordResetRepository is the concrete implementation of PasswordResetRepository interface
type passwordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository creates a new instance of passwordResetRepository
func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{
		db: db,
	}
}

func (r *passwordResetRepository) CreateToken(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *passwordResetRepository) GetToken(token string) (*models.PasswordResetToken, error) {
	var prt models.PasswordResetToken
	if err := r.db.First(&prt, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return &prt, nil
}

func (r *passwordResetRepository) DeleteToken(id uuid.UUID) error {
	return r.db.Delete(&models.PasswordResetToken{}, "id = ?", id).Error
}
