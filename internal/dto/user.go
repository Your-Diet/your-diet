package dto

import "time"

// RegisterUserRequest defines the expected request body for user registration.
// Email and Password validation is handled by the use case.
type RegisterUserRequest struct {
	Email    string  `json:"email" binding:"required"`
	Password string  `json:"password" binding:"required"`
}

// RegisterUserResponse defines the response for a successful user registration.
type RegisterUserResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

// LoginRequest defines the expected request body for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse defines the response for a successful user login.
type LoginResponse struct {
	Token        string    `json:"token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	UserType     string    `json:"user_type"`
	Permissions  []string  `json:"permissions"`
}
