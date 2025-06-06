package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the system
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Type     string             `bson:"type" json:"type"`
	Age      int                `bson:"age" json:"age"`
	Weight   float64            `bson:"weight" json:"weight"`
	Height   float64            `bson:"height" json:"height"`
	Gender   string             `bson:"gender" json:"gender"`
	Goal     string             `bson:"goal" json:"goal"`
	Activity string             `bson:"activity" json:"activity"`
	Calories int                `bson:"calories" json:"calories"`
	Protein  int                `bson:"protein" json:"protein"`
	Carbs    int                `bson:"carbs" json:"carbs"`
	Fat      int                `bson:"fat" json:"fat"`
}
