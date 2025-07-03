package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
)

type NotificationHandler struct {
	hub *usecase.NotificationHub
}

func NewNotificationHandler(hub *usecase.NotificationHub) *NotificationHandler {
	return &NotificationHandler{
		hub: hub,
	}
}

func (h *NotificationHandler) Handle(c *gin.Context) {
	var notification usecase.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Notify(notification)
	c.Status(http.StatusNoContent)
}
