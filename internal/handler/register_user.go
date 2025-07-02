package handler

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/dto"
	"github.com/victorgiudicissi/your-diet/internal/entity"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

var (
	// Regex for email validation: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	ErrInvalidEmailFormat     = errors.New("invalid email format")
	ErrPasswordTooShort       = errors.New("password must be at least 8 characters long")
	ErrPasswordMissingSpecial = errors.New("password must contain at least one special character")
)

// RegisterUserHandler handles HTTP requests related to users.
type RegisterUserHandler struct {
	createUserUseCase usecase.CreateUser
}

func NewRegisterUserHandler(createUserUC usecase.CreateUser) *RegisterUserHandler {
	return &RegisterUserHandler{
		createUserUseCase: createUserUC,
	}
}

// Handle handles the HTTP request for user registration using Gin.
func (h *RegisterUserHandler) Handle(c *gin.Context) {
	var req dto.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[RegisterUserHandler] Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, dto.NewError("something went wrong binding request data", err.Error()))
		return
	}
	// Validate Email
	if !emailRegex.MatchString(req.Email) {
		log.Printf("[RegisterUserHandler] Invalid email format: %s", req.Email)
		c.JSON(http.StatusBadRequest, dto.NewError("something went wrong validating request data", ErrInvalidEmailFormat.Error()))
		return
	}

	passworValidationErrs := validatePassword(req.Password)
	if len(passworValidationErrs) > 0 {
		log.Printf("[RegisterUserHandler] Password validation error: %s", passworValidationErrs[0])
		c.JSON(http.StatusBadRequest, dto.NewError("something went wrong validating request data", passworValidationErrs[0]))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[RegisterUserHandler] Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, dto.NewError("something went wrong hashing password", "failed to hash password"))
		return
	}

	userType := "DEFAULT"

	if req.IsNutritionist {
		userType = "NUTRITIONIST"
	}

	user := &entity.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Type:     userType,
		Age:      req.Age,
		Gender:   req.Gender,
	}

	err = h.createUserUseCase.Execute(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyExists) {
			log.Printf("[RegisterUserHandler] Email already exists: %s", req.Email)
			c.JSON(http.StatusConflict, dto.NewError("email already exists", "a user with this email already exists"))
			return
		}
		log.Printf("[RegisterUserHandler] Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, dto.NewError("something went wrong creating user", "failed to register user"))
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterUserResponse{
		Message: "user registered successfully",
	})
}

func validatePassword(password string) []string {
	var errors []string

	if len(password) < 8 || len(password) > 12 {
		errors = append(errors, "the password must be between 8 and 12 characters long")
	}

	letterCount := 0
	for _, r := range password {
		if unicode.IsLetter(r) {
			letterCount++
		}
	}
	if letterCount < 2 {
		errors = append(errors, "the password must contain at least 2 letters")
	}

	specialCharRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	if !specialCharRegex.MatchString(password) {
		errors = append(errors, "the password must contain at least one special character")
	}

	return errors
}
