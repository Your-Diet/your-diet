package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/dto"
	"github.com/victorgiudicissi/your-diet/internal/entity"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
)

type LoginHandler struct {
	loginUseCase usecase.LoginUseCase
}

func NewLoginHandler(loginUC usecase.LoginUseCase) *LoginHandler {
	return &LoginHandler{
		loginUseCase: loginUC,
	}
}

func (h *LoginHandler) HandleLogin(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewError("error", err.Error()))
		return
	}

	input := &entity.LoginUseCaseInput{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.loginUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, dto.NewError("error", err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.NewError("error", "login failed"))
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token:       result.Token,
		ExpiresAt:   result.ExpiresAt,
	})
}
