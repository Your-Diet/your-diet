package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/dto"
	"github.com/victorgiudicissi/your-diet/internal/middleware"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
)

// ListDietsHandler lida com as requisições de listagem de dietas
type ListDietsHandler struct {
	listDietsUseCase usecase.ListDiets
}

// NewListDietsHandler cria uma nova instância de ListDietsHandler
func NewListDietsHandler(listDietsUseCase usecase.ListDiets) *ListDietsHandler {
	return &ListDietsHandler{
		listDietsUseCase: listDietsUseCase,
	}
}

// Handle lida com a requisição de listagem de dietas
func (h *ListDietsHandler) Handle(c *gin.Context) {
	// Tenta obter as claims do contexto do Gin primeiro
	claimsValue, exists := c.Get(string(middleware.TokenContextKey))
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.NewError("something went wrong getting user claims", "token not found"))
		return
	}

	// Executar o caso de uso
	output, err := h.listDietsUseCase.Execute(c.Request.Context(), claimsValue.(*middleware.Claims).Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewError("something went wrong listing diets", err.Error()))
		return
	}

	c.JSON(http.StatusOK, output.Diets)
}
