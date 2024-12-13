// middleware/auth_middleware.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ABDULS21985/test-portal/services"
	"github.com/ABDULS21985/test-portal/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type key string

const (
	UserContext key = "user"
)

// Claims defines the structure of JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.StandardClaims
}

type AuthMiddleware struct {
	authService services.AuthService
}

func NewAuthMiddleware(authService services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// NewAuthMiddleware creates a new instance of AuthMiddleware with the given JWT secret
// NewJWTMiddleware creates a new instance of JWT middleware
func NewJWTMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserContext, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth ensures the user is authenticated
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header is missing")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := m.authService.ValidateToken(token)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireRole ensures the user has a specific role
func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header is missing")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := m.authService.GetClaimsFromToken(token)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			if claims["role"] != role {
				utils.RespondWithError(w, http.StatusForbidden, "You do not have permission to access this resource")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
