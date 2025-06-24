package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/dto"
	"github.com/victorgiudicissi/your-diet/internal/middleware"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
)

type CreateDietHandler struct {
	createDietUseCase usecase.CreateDiet
}

func NewCreateDietHandler(createDietUseCase usecase.CreateDiet) *CreateDietHandler {
	return &CreateDietHandler{
		createDietUseCase: createDietUseCase,
	}
}

func (h *CreateDietHandler) Handle(c *gin.Context) {
	var req dto.DietRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateDietHandler] Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, dto.NewError("something went wrong binding request data", err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		log.Printf("[CreateDietHandler] Validation error: %v", err)
		c.JSON(http.StatusBadRequest, dto.NewError("something went wrong validating request data", err.Error()))
		return
	}

	claimsValue, exists := c.Get(string(middleware.TokenContextKey))
	if !exists {
		log.Printf("[CreateDietHandler] Missing user claims in context")
		c.JSON(http.StatusUnauthorized, dto.NewError("something went wrong getting user claims", "usuário não autenticado"))
		return
	}

	diet, err := dto.ConvertToDiet(claimsValue.(*middleware.Claims).UserID, &req)
	if err != nil {
		log.Printf("[CreateDietHandler] Failed to convert to diet: %v", err)
		c.JSON(http.StatusBadRequest, dto.NewError("something went wrong creating diet", "invalid ingredients: "+err.Error()))
		return
	}

	if err := h.createDietUseCase.Execute(c.Request.Context(), diet); err != nil {
		log.Printf("[CreateDietHandler] Failed to execute use case: %v", err)
		c.JSON(http.StatusInternalServerError, dto.NewError("something went wrong creating diet", "failed to create diet request: "+err.Error()))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
