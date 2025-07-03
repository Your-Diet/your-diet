package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/victorgiudicissi/your-diet/internal/entity"
)

// UpdateDietUseCaseInput contém os dados de entrada para atualização de dieta
type UpdateDietUseCaseInput struct {
	DietID         string
	UserEmail      string
	DietName       *string
	DurationInDays *uint32
	Status         *string
	Meals          []entity.Meal
	Observations   *string
}

// UpdateDietUseCase define a interface para o caso de uso de atualização de dieta
type UpdateDietUseCase interface {
	Execute(ctx context.Context, dietID string, newDiet *entity.Diet) (*entity.Diet, error)
}

type updateDietUseCase struct {
	dietRepo DietRepository
}

// NewUpdateDiet cria uma nova instância de UpdateDietUseCase
func NewUpdateDiet(dietRepo DietRepository) UpdateDietUseCase {
	return &updateDietUseCase{
		dietRepo: dietRepo,
	}
}

func (uc *updateDietUseCase) Execute(ctx context.Context, dietID string, newDiet *entity.Diet) (*entity.Diet, error) {
	diet, err := uc.dietRepo.GetDietByID(ctx, dietID)
	if err != nil {
		return nil, err
	}

	if newDiet.CreatedBy != diet.CreatedBy {
		return nil, errors.New("usuário não autorizado a atualizar esta dieta")
	}

	if newDiet.DietName != "" && newDiet.DietName != diet.DietName {
		diet.DietName = newDiet.DietName
	}

	if newDiet.DurationInDays != 0 && newDiet.DurationInDays != diet.DurationInDays {
		diet.DurationInDays = newDiet.DurationInDays
	}

	if newDiet.Status != "" && newDiet.Status != diet.Status {
		diet.Status = newDiet.Status
	}

	if len(newDiet.Meals) > 0 {
		diet.Meals = newDiet.Meals
	}

	if newDiet.Observations != "" && newDiet.Observations != diet.Observations {
		diet.Observations = newDiet.Observations
	}

	diet.UpdatedAt = time.Now()

	if err := uc.dietRepo.UpdateDiet(ctx, diet); err != nil {
		return nil, err
	}

	return diet, nil
}
