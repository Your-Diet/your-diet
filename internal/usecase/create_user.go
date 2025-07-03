package usecase

import (
	"context"
	"errors"

	"github.com/victorgiudicissi/your-diet/internal/entity"
)

var ErrEmailAlreadyExists = errors.New("a user with this email already exists")

// CreateUserUseCase handles the logic for creating a new user.
type createUserUseCase struct {
	userRepo UserRepository
}

type CreateUser interface {
	Execute(ctx context.Context, user *entity.User) error
}

// NewCreateUser creates a new instance of CreateUserUseCase.
func NewCreateUser(userRepo UserRepository) CreateUser {
	return &createUserUseCase{userRepo: userRepo}
}

// Execute creates a new user.
func (uc *createUserUseCase) Execute(ctx context.Context, user *entity.User) error {
	existing, err := uc.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrEmailAlreadyExists
	}

	_, err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
