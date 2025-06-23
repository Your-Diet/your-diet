package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the system
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Type     string             `bson:"type" json:"type"`
	Age      int                `bson:"age" json:"age"`
	Gender   string             `bson:"gender" json:"gender"`
}
