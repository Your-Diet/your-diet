package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/dto"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
)

// LoginHandler handles HTTP requests related to user login.
type LoginHandler struct {
	loginUseCase *usecase.LoginUseCase
}

func NewLoginHandler(loginUC *usecase.LoginUseCase) *LoginHandler {
	return &LoginHandler{
		loginUseCase: loginUC,
	}
}

// HandleLogin handles the HTTP request for user login using Gin.
func (h *LoginHandler) HandleLogin(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		log.Printf("Error binding JSON for login: %v", err)
		return
	}

	input := usecase.LoginUseCaseInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.loginUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		log.Printf("Error during login: %v", err)
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token:       output.Token,
		ExpiresAt:   output.ExpiresAt,
		UserID:      output.UserID,
		Email:       output.Email,
		UserType:    output.UserType,
		Permissions: output.Permissions,
	})
}
