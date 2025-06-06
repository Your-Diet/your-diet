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
)

// DietRepository define a interface para o repositório de dietas
type DietRepository interface {
	CreateDiet(ctx context.Context, diet *entity.Diet) error
	GetDietByID(ctx context.Context, id string) (*entity.Diet, error)
	FindByUserEmail(ctx context.Context, userEmail string) ([]*entity.Diet, error)
	UpdateDiet(ctx context.Context, diet *entity.Diet) error
	Close() error
}

const ( 
	dietCollectionName = "diets"
)

type dietRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewDietRepository(uri, database string) (DietRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &dietRepository{
		client:     client,
		database:   database,
		collection: dietCollectionName,
	}, nil
}

func (r *dietRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return r.client.Disconnect(ctx)
}

func (r *dietRepository) CreateDiet(ctx context.Context, diet *entity.Diet) error {
	collection := r.client.Database(r.database).Collection(r.collection)

	diet.CreatedAt = time.Now()
	diet.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, diet)
	if err != nil {
		return err
	}

	return nil
}

func (r *dietRepository) GetDietByID(ctx context.Context, id string) (*entity.Diet, error) {
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

// FindByUserEmail retorna todas as dietas de um usuário
func (r *dietRepository) FindByUserEmail(ctx context.Context, userEmail string) ([]*entity.Diet, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	cursor, err := collection.Find(ctx, bson.M{"user_email": userEmail})
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
func (r *dietRepository) UpdateDiet(ctx context.Context, diet *entity.Diet) error {
	collection := r.client.Database(r.database).Collection(r.collection)

	diet.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":              diet.DietName,
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
