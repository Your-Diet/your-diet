package usecase

import (
	"context"

	"github.com/victorgiudicissi/your-diet/internal/entity"
)

type CreateDietUseCase interface {
	Execute(ctx context.Context, diet *entity.Diet) error
}

type createDietUseCase struct {
	dietRepo DietRepository
}

func NewCreateDietUseCase(dietRepo DietRepository) CreateDietUseCase {
	return &createDietUseCase{
		dietRepo: dietRepo,
	}
}

func (uc *createDietUseCase) Execute(ctx context.Context, diet *entity.Diet) error {
	diet.Status = "pending"

	if err := uc.dietRepo.CreateDiet(ctx, diet); err != nil {
		return err
	}

	return nil
}
