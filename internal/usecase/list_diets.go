package usecase

import (
	"context"

	"github.com/victorgiudicissi/your-diet/internal/dto"
	"github.com/victorgiudicissi/your-diet/internal/repository"
)

type listDietsUseCase struct {
	dietRepo repository.DietRepository
}

type ListDiets interface{
	Execute(ctx context.Context, userEmail string) (*dto.ListDietsUseCaseOutput, error)
}

func NewListDietsUseCase(dietRepo repository.DietRepository) ListDiets {
	return &listDietsUseCase{dietRepo: dietRepo}
}

func (uc *listDietsUseCase) Execute(ctx context.Context, userEmail string) (*dto.ListDietsUseCaseOutput, error) {
	diets, err := uc.dietRepo.FindByUserEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}
	
	return dto.NewListDietsUseCaseOutput(diets), nil
}
