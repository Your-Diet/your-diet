package dto

import "time"

// RegisterUserRequest defines the expected request body for user registration.
// Email and Password validation is handled by the use case.
type RegisterUserRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required"`
	Age            int    `json:"age" binding:"required,min=1"`
	Gender         string `json:"gender" binding:"required,oneof=male female other"`
	IsNutritionist bool   `json:"is_nutritionist"`
}

// RegisterUserResponse defines the response for a successful user registration.
type RegisterUserResponse struct {
	Message string `json:"message"`
}

// LoginRequest defines the expected request body for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse defines the response for a successful user login.
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
