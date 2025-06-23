package usecase

import (
	"context"

	"github.com/victorgiudicissi/your-diet/internal/dto"
)

type DietFilter struct {
	UserEmail *string
	CreatedBy *string
}

type listDietsUseCase struct {
	dietRepo DietRepository
	userRepo UserRepository
}

type ListDiets interface{
	Execute(ctx context.Context, input *dto.ListDietsInput) (*dto.ListDietsUseCaseOutput, error)
}

func NewListDietsUseCase(dietRepo DietRepository, userRepo UserRepository) ListDiets {
	return &listDietsUseCase{
		dietRepo: dietRepo,
		userRepo: userRepo,
	}
}

func (uc *listDietsUseCase) Execute(ctx context.Context, input *dto.ListDietsInput) (*dto.ListDietsUseCaseOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	filter := &DietFilter{}

	filter.UserEmail = &user.Email
	if input.CreatedBySearch && user.Type == "NUTRITIONIST" {
		filter.CreatedBy = &input.UserID
		if input.UserEmail != "" {
			filter.UserEmail = &input.UserEmail
		}
	}

	diets, err := uc.dietRepo.FindDiets(ctx, filter)
	if err != nil {
		return nil, err
	}

	return dto.NewListDietsUseCaseOutput(diets), nil
}
