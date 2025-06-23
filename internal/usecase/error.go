package usecase

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized user")
	ErrUserNotFound = errors.New("user not found")
)
