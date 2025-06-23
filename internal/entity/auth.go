package entity

import "time"

type LoginUseCaseInput struct {
	Email    string
	Password string
}

type LoginUseCaseOutput struct {
	Token       string
	ExpiresAt   time.Time
}