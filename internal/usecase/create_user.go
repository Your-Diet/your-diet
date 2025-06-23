package usecase

import (
	"context"

	"github.com/victorgiudicissi/your-diet/internal/entity"
)

// CreateUserUseCase handles the logic for creating a new user.
type createUserUseCase struct {
	userRepo UserRepository
}

type CreateUser interface {
	Execute(ctx context.Context, user *entity.User) error
}

// NewCreateUserUseCase creates a new instance of CreateUserUseCase.
func NewCreateUserUseCase(userRepo UserRepository) CreateUser {
	return &createUserUseCase{userRepo: userRepo}
}

// Execute creates a new user.
func (uc *createUserUseCase) Execute(ctx context.Context, user *entity.User) error {
	_, err := uc.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
