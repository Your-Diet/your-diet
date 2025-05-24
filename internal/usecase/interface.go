package usecase

import (
	"context"

	"github.com/victorgiudicissi/your-diet/internal/entity"
)

type (
	DietRepository interface {
		CreateDiet(ctx context.Context, diet *entity.Diet) error
		GetDietByID(ctx context.Context, id string) (*entity.Diet, error)
	}
)
