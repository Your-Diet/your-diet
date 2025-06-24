package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/victorgiudicissi/your-diet/internal/entity"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
	"github.com/victorgiudicissi/your-diet/internal/utils"
)

const (
	dietCollectionName = "diets"
)

type DietRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewDietRepository(cfg *utils.EnvConfig) (*DietRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURL))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &DietRepository{
		client:     client,
		database:   cfg.DBName,
		collection: dietCollectionName,
	}, nil
}

func (r *DietRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return r.client.Disconnect(ctx)
}

func (r *DietRepository) CreateDiet(ctx context.Context, diet *entity.Diet) error {
	collection := r.client.Database(r.database).Collection(r.collection)

	diet.CreatedAt = time.Now()
	diet.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, diet)
	if err != nil {
		return err
	}

	return nil
}

func (r *DietRepository) GetDietByID(ctx context.Context, id string) (*entity.Diet, error) {
	collection := r.client.Database(r.database).Collection(r.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var diet entity.Diet
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&diet)
	if err != nil {
		return nil, err
	}

	return &diet, nil
}

// FindDiets retorna as dietas com base nos filtros fornecidos
func (r *DietRepository) FindDiets(ctx context.Context, filter *usecase.DietFilter) ([]*entity.Diet, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	// Construir o filtro do MongoDB dinamicamente com base nos campos fornecidos
	mongoFilter := bson.M{}

	if filter.UserEmail != nil {
		mongoFilter["user_email"] = *filter.UserEmail
	}

	if filter.CreatedBy != nil {
		mongoFilter["created_by"] = *filter.CreatedBy
	}

	cursor, err := collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var diets []*entity.Diet
	if err = cursor.All(ctx, &diets); err != nil {
		return nil, err
	}

	return diets, nil
}

// UpdateDiet atualiza uma dieta existente
func (r *DietRepository) UpdateDiet(ctx context.Context, diet *entity.Diet) error {
	collection := r.client.Database(r.database).Collection(r.collection)

	diet.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":             diet.DietName,
			"duration_in_days": diet.DurationInDays,
			"status":           diet.Status,
			"meals":            diet.Meals,
			"observations":     diet.Observations,
			"updated_at":       diet.UpdatedAt,
		},
	}

	objID, err := primitive.ObjectIDFromHex(diet.ID)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objID, "user_email": diet.UserEmail}, // Garante que só o dono pode atualizar
		update,
	)

	return err
}
