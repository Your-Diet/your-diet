package handler

import (
	"log"
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
		log.Printf("[ListDietsHandler] Missing user claims in context")
		c.JSON(http.StatusUnauthorized, dto.NewError("something went wrong getting user claims", "usuário não autenticado"))
	}

	// Obter os parâmetros da query string
	userEmail := c.Query("userEmail")
	createdBySearch := c.Query("createdBySearch") == "true"
	userID := claimsValue.(*middleware.Claims).UserID

	// Criar o input para o caso de uso
	input := &dto.ListDietsInput{
		UserEmail:       userEmail,
		CreatedBySearch: createdBySearch,
		UserID:          userID,
	}

	// Executar o caso de uso
	output, err := h.listDietsUseCase.Execute(c.Request.Context(), input)

	if err != nil {
		log.Printf("[ListDietsHandler] Failed to list diets: %v", err)
		c.JSON(http.StatusInternalServerError, dto.NewError("something went wrong listing diets", err.Error()))
		return
	}

	c.JSON(http.StatusOK, output.Diets)
}
