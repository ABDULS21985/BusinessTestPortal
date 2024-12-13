// services/auth_service.go
package services

import (
	"errors"

	"github.com/ABDULS21985/test-portal/models"
	"github.com/ABDULS21985/test-portal/repositories"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the contract for authentication-related business logic
type AuthService interface {
	AuthenticateUser(email, password string) (*models.User, string, error)
	GenerateToken(userID uuid.UUID, role string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GetClaimsFromToken(token string) (jwt.MapClaims, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret []byte
}

// NewAuthService creates a new instance of authService
func NewAuthService(userRepo repositories.UserRepository, jwtSecret []byte) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.StandardClaims
}

func (s *authService) AuthenticateUser(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", err
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}

// GenerateToken generates a JWT for a given user ID and role
func (a *authService) GenerateToken(userID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token valid for 72 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

// ValidateToken validates a given JWT token
func (a *authService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return a.jwtSecret, nil
	})
}

// GetClaimsFromToken extracts claims from a validated token
func (a *authService) GetClaimsFromToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := a.ValidateToken(token)
	if err != nil || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("unable to parse claims")
	}

	return claims, nil
}
