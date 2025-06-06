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

var ErrUserNotFound = errors.New("user not found")	

// MongoUserRepository implements the usecase.UserRepository interface using MongoDB.
type MongoUserRepository struct {
	client     *mongo.Client
	dbName     string
	collection *mongo.Collection
}

// NewMongoUserRepository creates a new MongoUserRepository.
func NewMongoUserRepository(mongoURI, dbName string) (*MongoUserRepository, error) {
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

	collection := client.Database(dbName).Collection(userCollectionName)
	return &MongoUserRepository{
		client:     client,
		dbName:     dbName,
		collection: collection,
	}, nil
}

// Create inserts a new user into the database.
// It returns the ID of the newly created user or an error.
func (r *MongoUserRepository) Create(ctx context.Context, user *entity.User) (string, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	// Return the string representation of the inserted ID
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Disconnect closes the MongoDB connection.
func (r *MongoUserRepository) Disconnect(ctx context.Context) error {
	if r.client != nil {
		return r.client.Disconnect(ctx)
	}
	return nil
}
