package dto

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// DietRequest represents the request body for creating a new diet
type DietRequest struct {
	UserEmail      string `json:"user_email" validate:"required,email"`
	DietName       string `json:"name" validate:"required,min=3,max=100"`
	DurationInDays uint32 `json:"duration_in_days" validate:"required,min=1"`
}

// ValidationError represents a custom validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return e.Message
}

func fieldNameToHumanReadable(field string) string {
	switch field {
	case "UserEmail":
		return "email"
	case "DietName":
		return "diet name"
	case "DurationInDays":
		return "duration"
	default:
		var result strings.Builder
		for i, r := range field {
			if i > 0 && 'A' <= r && r <= 'Z' {
				result.WriteRune(' ')
			}
			result.WriteRune(r)
		}
		return strings.ToLower(result.String())
	}
}

func (d *DietRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(d)

	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := fieldNameToHumanReadable(fieldError.Field())
			switch fieldError.Tag() {
			case "required":
				return &ValidationError{
					Field:   fieldError.Field(),
					Message: fmt.Sprintf("The %s field is required", fieldName),
				}
			case "email":
				return &ValidationError{
					Field:   fieldError.Field(),
					Message: "Please provide a valid email address",
				}
			case "min":
				if fieldError.Field() == "DurationInDays" {
					return &ValidationError{
						Field:   fieldError.Field(),
						Message: "The duration must be at least 1 day",
					}
				}
				return &ValidationError{
					Field:   fieldError.Field(),
					Message: fmt.Sprintf("The %s must be at least %s characters", fieldName, fieldError.Param()),
				}
			case "max":
				return &ValidationError{
					Field:   fieldError.Field(),
					Message: fmt.Sprintf("The %s must not exceed %s characters", fieldName, fieldError.Param()),
				}
			}
		}
	}

	return &ValidationError{
		Field:   "",
		Message: "Invalid request data",
	}
}
