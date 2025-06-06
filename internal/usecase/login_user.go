package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/victorgiudicissi/your-diet/internal/repository"
)

// TODO: Move this to a secure configuration (e.g., environment variable)
var jwtSecretKey = []byte("your-very-secret-jwt-key-here")

const (
	TokenTypeDefault     = "DEFAULT"
	TokenTypeNutritionist = "NUTRITIONIST"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotActive      = errors.New("user account is not active") // Placeholder for future use
)

// LoginUseCaseInput represents the input data for logging in a user.
type LoginUseCaseInput struct {
	Email    string
	Password string
}

// LoginUseCaseOutput represents the output data after a successful login.
type LoginUseCaseOutput struct {
	Token       string
	ExpiresAt   time.Time
	UserID      string
	Email       string
	UserType    string
	Permissions []string
}

// Claims defines the JWT claims structure.
type Claims struct {
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// LoginUseCase handles the logic for user login.
type LoginUseCase struct {
	userRepo UserRepository
}

// NewLoginUseCase creates a new instance of LoginUseCase.
func NewLoginUseCase(userRepo UserRepository) *LoginUseCase {
	return &LoginUseCase{userRepo: userRepo}
}

// Execute performs the login operation.
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginUseCaseInput) (*LoginUseCaseOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err // Other repository error
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		// Handles incorrect password and other bcrypt errors
		return nil, ErrInvalidCredentials
	}

	// Determine permissions based on user type
	var permissions []string
	switch user.Type {
	case TokenTypeDefault:
		permissions = []string{"list_diet"}
	case TokenTypeNutritionist:
		permissions = []string{"list_diet", "create_diet", "upload_file", "update_diet"}
	default:
		// Default to no permissions or minimal permissions if type is unknown or not set
		permissions = []string{}
	}

	// Create JWT token
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		UserID:      user.ID.Hex(),
		Email:       user.Email,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return nil, err // Internal server error (failed to sign token)
	}

	return &LoginUseCaseOutput{
		Token:       tokenString,
		ExpiresAt:   expirationTime,
		UserID:      user.ID.Hex(),
		Email:       user.Email,
		UserType:    user.Type,
		Permissions: permissions,
	}, nil
}
