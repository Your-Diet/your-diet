package usecase

import (
	"context"

	"github.com/victorgiudicissi/your-diet/internal/entity"
)


// CreateUserUseCaseOutput represents the output data after creating a user.
type CreateUserUseCaseOutput struct {
	UserID string
}

// CreateUserUseCase handles the logic for creating a new user.
type CreateUserUseCase struct {
	userRepo UserRepository
}

// NewCreateUserUseCase creates a new instance of CreateUserUseCase.
func NewCreateUserUseCase(userRepo UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo}
}

// Execute creates a new user.
func (uc *CreateUserUseCase) Execute(ctx context.Context, user *entity.User) (*CreateUserUseCaseOutput, error) {
	user.Type = "DEFAULT"
	userID, err := uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &CreateUserUseCaseOutput{UserID: userID}, nil
}
