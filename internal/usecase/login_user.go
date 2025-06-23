package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/victorgiudicissi/your-diet/internal/constants"
	"github.com/victorgiudicissi/your-diet/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

var JWTSecretKey = []byte("your-very-secret-jwt-key-here")

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotActive      = errors.New("user account is not active")
)

type LoginUseCase interface {
	Execute(ctx context.Context, input *entity.LoginUseCaseInput) (*entity.LoginUseCaseOutput, error)
}

type Claims struct {
	UserID      string   `json:"user_id"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type loginUseCase struct {
	userRepo UserRepository
}

func NewLoginUseCase(userRepo UserRepository) LoginUseCase {
	return &loginUseCase{userRepo: userRepo}
}

func (uc *loginUseCase) Execute(ctx context.Context, input *entity.LoginUseCaseInput) (*entity.LoginUseCaseOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	permissions := constants.GetPermissionsByUserType(user.Type)

	expirationTime := time.Now().Add(1000 * time.Minute)
	claims := &Claims{
		UserID:      user.ID.Hex(),
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return nil, err // Internal server error (failed to sign token)
	}

	return &entity.LoginUseCaseOutput{
		Token:       tokenString,
		ExpiresAt:   expirationTime,
	}, nil
}
