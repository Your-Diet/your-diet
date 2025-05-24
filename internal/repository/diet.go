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
)

type DietRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewDietRepository(uri, database, collection string) (usecase.DietRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &DietRepository{
		client:     client,
		database:   database,
		collection: collection,
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
