package dto

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/victorgiudicissi/your-diet/internal/entity"
)

// IngredientRequest representa um ingrediente na requisição
type IngredientRequest struct {
	Description string              `json:"description" validate:"required,min=1"`
	Quantity    float64             `json:"quantity" validate:"required,min=0"`
	Unit        string              `json:"unit" validate:"required,oneof=ml g l kg mg un fatia(s)"`
	Substitutes []IngredientRequest `json:"substitutes"`
}

// MealRequest representa uma refeição na requisição
type MealRequest struct {
	Name        string              `json:"name" validate:"required,min=3,max=100"`
	Description string              `json:"description"`
	TimeOfDay   string              `json:"time_of_day" validate:"required,oneof=café_da_manhã almoço jantar lanche ceia"`
	Ingredients []IngredientRequest `json:"ingredients" validate:"required,min=1,dive"`
}

// DietRequest represents the request body for creating a new diet
type DietRequest struct {
	UserEmail      string        `json:"user_email" validate:"required,email"`
	DietName       string        `json:"name" validate:"required,min=3,max=100"`
	DurationInDays uint32        `json:"duration_in_days" validate:"required,min=1"`
	Meals          []MealRequest `json:"meals" validate:"required,min=1,dive"`
	Observations   string        `json:"observations"`
}

func ConvertToDiet(createdBy string, req *DietRequest) (*entity.Diet, error) {
	now := time.Now()

	meals := make([]entity.Meal, 0, len(req.Meals))
	for _, mealReq := range req.Meals {
		meal, err := ConvertToMeal(&mealReq)
		if err != nil {
			return nil, err
		}
		meals = append(meals, *meal)
	}

	return &entity.Diet{
		UserEmail:      req.UserEmail,
		DietName:       req.DietName,
		DurationInDays: req.DurationInDays,
		Status:         string(entity.Enabled),
		Meals:          meals,
		Observations:   req.Observations,
		CreatedBy:      createdBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func ConvertToMeal(req *MealRequest) (*entity.Meal, error) {
	ingredients := make([]entity.Ingredient, 0, len(req.Ingredients))
	for _, ingReq := range req.Ingredients {
		ingredient, err := ConvertToIngredient(&ingReq)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, *ingredient)
	}

	return &entity.Meal{
		Name:        req.Name,
		Description: req.Description,
		TimeOfDay:   req.TimeOfDay,
		Ingredients: ingredients,
	}, nil
}

func ConvertToIngredient(req *IngredientRequest) (*entity.Ingredient, error) {
	substitutes := make([]entity.Ingredient, 0, len(req.Substitutes))
	for _, subReq := range req.Substitutes {
		subIngredient, err := ConvertToIngredient(&subReq)
		if err != nil {
			return nil, err
		}
		substitutes = append(substitutes, *subIngredient)
	}

	return &entity.Ingredient{
		Description: req.Description,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		Substitutes: substitutes,
	}, nil
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
		return "nome da dieta"
	case "DurationInDays":
		return "duração em dias"
	case "Meals":
		return "refeições"
	case "Name":
		return "nome"
	case "TimeOfDay":
		return "período do dia"
	case "Ingredients":
		return "ingredientes"
	case "Description":
		return "descrição"
	case "Quantity":
		return "quantidade"
	case "Unit":
		return "unidade de medida"
	case "Observations":
		return "observações"
	default:
		return field
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

type DietResponse struct {
	ID             string   
	UserEmail      string   
	DietName       string   
	DurationInDays uint32   
	Status         string   
	Meals         []MealResponse    
	Observations   string    
	CreatedBy      string    
	CreatedAt      time.Time 
	UpdatedAt      time.Time 
}

type MealResponse struct {
	Name        string      
	Description string      
	TimeOfDay   string      
	Ingredients []IngredientResponse 
}

type IngredientResponse struct {
	Description string      
	Quantity    float64     
	Unit        string      
	Substitutes []IngredientResponse 
}

type ListDietsUseCaseOutput struct {
	Diets []*DietResponse
}

func NewListDietsUseCaseOutput(diets []*entity.Diet) *ListDietsUseCaseOutput {
	var dietsResponse []*DietResponse
	for _, diet := range diets {
		dietsResponse = append(dietsResponse, &DietResponse{
			ID:             diet.ID,
			UserEmail:      diet.UserEmail,
			DietName:       diet.DietName,
			DurationInDays: diet.DurationInDays,
			Status:         diet.Status,
			Meals:          convertMealsToMealResponse(diet.Meals),
			Observations:   diet.Observations,
			CreatedBy:      diet.CreatedBy,
			CreatedAt:      diet.CreatedAt,
			UpdatedAt:      diet.UpdatedAt,
		})
	}

	return &ListDietsUseCaseOutput{
		Diets: dietsResponse,
	}
}

func convertMealsToMealResponse(meals []entity.Meal) []MealResponse {
	var mealResponses []MealResponse
	for _, meal := range meals {
		mealResponses = append(mealResponses, MealResponse{
			Name:        meal.Name,
			Description: meal.Description,
			TimeOfDay:   meal.TimeOfDay,
			Ingredients: convertIngredientsToIngredientResponse(meal.Ingredients),
		})
	}
	return mealResponses
}

func convertIngredientsToIngredientResponse(ingredients []entity.Ingredient) []IngredientResponse {
	var ingredientResponses []IngredientResponse
	for _, ingredient := range ingredients {
		ingredientResponses = append(ingredientResponses, IngredientResponse{
			Description: ingredient.Description,
			Quantity:    ingredient.Quantity,
			Unit:        ingredient.Unit,
			Substitutes: convertIngredientsToIngredientResponse(ingredient.Substitutes),
		})
	}
	return ingredientResponses
}
