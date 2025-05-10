package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"dnd-combat/config"
	"dnd-combat/internal/models"
)

// Error definitions
var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Service handles authentication business logic
type Service struct {
	repo   *Repository
	config *config.Config
}

// NewService creates a new auth service
func NewService(repo *Repository, config *config.Config) *Service {
	return &Service{
		repo:   repo,
		config: config,
	}
}

// Register creates a new user account
func (s *Service) Register(user *models.User, password string) error {
	// Check if username is already taken
	existingUser, err := s.repo.GetByUsername(user.Username)
	if err != nil {
		return fmt.Errorf("error checking for existing user: %w", err)
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	// Check if email is already taken
	existingUser, err = s.repo.GetByEmail(user.Email)
	if err != nil {
		return fmt.Errorf("error checking for existing email: %w", err)
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	user.PasswordHash = string(hashedPassword)

	// Create the user
	if err := s.repo.Create(user); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(username, password string) (*models.User, string, error) {
	// Fetch user by username
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching user: %w", err)
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("error generating token: %w", err)
	}

	return user, token, nil
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// generateToken creates a new JWT token for the given user
func (s *Service) generateToken(user *models.User) (string, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dnd-combat-api",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	signedToken, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
