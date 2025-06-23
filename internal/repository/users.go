package repository

import (
	"context"
	"errors"
	"time"

	"github.com/victorgiudicissi/your-diet/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userCollectionName = "users"
)

// UserRepository implements the usecase.UserRepository interface using MongoDB.
type UserRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

// NewMongoUserRepository creates a new MongoUserRepository.
func NewMongoUserRepository(mongoURI, dbName string) (*UserRepository, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &UserRepository{
		client:     client,
		database:   dbName,
		collection: userCollectionName,
	}, nil
}

// Create inserts a new user into the database.
// It returns the ID of the newly created user or an error.
func (r *UserRepository) Create(ctx context.Context, user *entity.User) (string, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	// Return the string representation of the inserted ID
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	var user entity.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	collection := r.client.Database(r.database).Collection(r.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user entity.User
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
