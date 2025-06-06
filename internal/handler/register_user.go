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
	createUserUseCase *usecase.CreateUserUseCase
}

func NewRegisterUserHandler(createUserUC *usecase.CreateUserUseCase) *RegisterUserHandler {
	return &RegisterUserHandler{
		createUserUseCase: createUserUC,
	}
}

// Handle handles the HTTP request for user registration using Gin.
func (h *RegisterUserHandler) Handle(c *gin.Context) {
	var req dto.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		log.Printf("Error binding JSON: %v", err)
		return
	}
	// Validate Email
	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidEmailFormat.Error()})
		return
	}

	passworValidationErrs := validatePassword(req.Password)
	if len(passworValidationErrs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": passworValidationErrs[0]})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := &entity.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	output, err := h.createUserUseCase.Execute(c.Request.Context(), user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterUserResponse{
		Message: "User registered successfully",
		UserID:  output.UserID,
	})
}

func validatePassword(password string) []string {
	var errors []string

	// 1. Verificar o tamanho da senha (entre 8 e 12 caracteres)
	if len(password) < 8 || len(password) > 12 {
		errors = append(errors, "A senha deve ter entre 8 e 12 caracteres.")
	}

	// 2. Verificar se há pelo menos 2 letras
	letterCount := 0
	for _, r := range password {
		if unicode.IsLetter(r) { // unicode.IsLetter() verifica letras de qualquer alfabeto
			letterCount++
		}
	}
	if letterCount < 2 {
		errors = append(errors, "A senha deve conter pelo menos 2 letras.")
	}

	// 3. Verificar se há pelo menos 1 caractere especial
	// Esta regex busca qualquer caractere que NÃO seja letra (a-z, A-Z) ou dígito (0-9).
	specialCharRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	if !specialCharRegex.MatchString(password) {
		errors = append(errors, "A senha deve conter pelo menos 1 caractere especial.")
	}

	return errors
}
