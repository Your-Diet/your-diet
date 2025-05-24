package entity

import "time"

type Diet struct {
	ID             string    `bson:"_id,omitempty" json:"id"`
	UserEmail      string    `bson:"user_email" json:"user_email"`
	DietName       string    `bson:"name" json:"name"`
	DurationInDays uint32    `bson:"duration_in_days" json:"duration_in_days"`
	Status         string    `bson:"status" json:"status"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}

func NewDietRequest(email, name string, duration uint32) *Diet {
	now := time.Now()
	return &Diet{
		UserEmail:      email,
		DietName:       name,
		DurationInDays: duration,
		Status:         "pending",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
