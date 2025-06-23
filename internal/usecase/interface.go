package usecase

import (
	"context"

	"github.com/victorgiudicissi/your-diet/internal/entity"
)

type (
	DietRepository interface {
		CreateDiet(ctx context.Context, diet *entity.Diet) error
		GetDietByID(ctx context.Context, id string) (*entity.Diet, error)
		FindDiets(ctx context.Context, filter *DietFilter) ([]*entity.Diet, error)
		UpdateDiet(ctx context.Context, diet *entity.Diet) error
	}

	UserRepository interface {
		Create(ctx context.Context, user *entity.User) (string, error)
		FindByEmail(ctx context.Context, email string) (*entity.User, error)
		FindByID(ctx context.Context, id string) (*entity.User, error)
	}
)
