package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/dto"
	"github.com/victorgiudicissi/your-diet/internal/entity"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
)

type CreateDietHandler struct {
	createDietUseCase usecase.CreateDietUseCase
}

func NewCreateDietHandler(createDietUseCase usecase.CreateDietUseCase) *CreateDietHandler {
	return &CreateDietHandler{
		createDietUseCase: createDietUseCase,
	}
}

func (h *CreateDietHandler) Handle(c *gin.Context) {
	var req dto.DietRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.createDietUseCase.Execute(c.Request.Context(), entity.NewDietRequest(
		req.UserEmail,
		req.DietName,
		req.DurationInDays,
	)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create diet request: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
