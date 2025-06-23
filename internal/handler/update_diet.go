package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/dto"
	"github.com/victorgiudicissi/your-diet/internal/middleware"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
)

type UpdateDietHandler struct {
	updateDietUseCase usecase.UpdateDietUseCase
}

func NewUpdateDietHandler(updateDietUseCase usecase.UpdateDietUseCase) *UpdateDietHandler {
	return &UpdateDietHandler{
		updateDietUseCase: updateDietUseCase,
	}
}

func (h *UpdateDietHandler) Handle(c *gin.Context) {
	// Obter o ID da dieta da URL
	dietID := c.Param("id")
	if dietID == "" {
		c.JSON(http.StatusBadRequest, dto.NewError("something went wrong updating diet", "ID da dieta é obrigatório"))
		return
	}

	// Obter o email do usuário do token JWT (já validado pelo middleware de autenticação)
	claimsValue, exists := c.Get(string(middleware.TokenContextKey))
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.NewError("something went wrong getting user claims", "usuário não autenticado"))
		return
	}

	// Fazer o bind do JSON para o DTO
	var req dto.DietRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewError("something went wrong updating diet", "dados inválidos: " + err.Error()))
		return
	}

	diet, err := dto.ConvertToDiet(claimsValue.(*middleware.Claims).UserID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewError("something went wrong updating diet", err.Error()))
		return
	}

	// Chamar o caso de uso
	updatedDiet, err := h.updateDietUseCase.Execute(c.Request.Context(), dietID, diet)
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "failed to update diet: " + err.Error()

		if errors.Is(err, usecase.ErrUnauthorized) {
			status = http.StatusForbidden
			errMsg = "you do not have permission to update this diet"
		}

		c.JSON(status, dto.NewError("something went wrong updating diet", errMsg))
		return
	}

	c.JSON(http.StatusOK, updatedDiet)
}
