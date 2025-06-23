package entity

import (
	"time"
)

type DietStatus string

const (
	Enabled  DietStatus = "ENABLED"
	Disabled DietStatus = "DISABLED"
)

type Diet struct {
	ID             string    `bson:"_id,omitempty" json:"id"`
	UserEmail      string    `bson:"user_email" json:"user_email"`
	DietName       string    `bson:"name" json:"name"`
	DurationInDays uint32    `bson:"duration_in_days" json:"duration_in_days"`
	Status         string    `bson:"status" json:"status"`
	Meals          []Meal    `bson:"meals" json:"meals"`
	Observations   string    `bson:"observations" json:"observations"`
	CreatedBy      string    `bson:"created_by" json:"created_by"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}

type Meal struct {
	Name        string       `bson:"name" json:"name"`
	Description string       `bson:"description" json:"description"`
	TimeOfDay   string       `bson:"time_of_day" json:"time_of_day"`
	Ingredients []Ingredient `bson:"ingredients" json:"ingredients"`
}

type Ingredient struct {
	Description string       `bson:"description" json:"description"`
	Quantity    float64      `bson:"quantity" json:"quantity"`
	Unit        string       `bson:"unit" json:"unit"`
	Substitutes []Ingredient `bson:"substitutes" json:"substitutes"`
}
