package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/middleware"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
)

type SSEHandler struct {
	hub *usecase.NotificationHub
}

func NewSSEHandler(hub *usecase.NotificationHub) *SSEHandler {
	return &SSEHandler{
		hub: hub,
	}
}

func (h *SSEHandler) Handle(c *gin.Context) {
	userID := c.GetString(string(middleware.UserIDContextKey))
	if userID == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	subscriber, err := h.hub.Subscribe(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-c.Writer.CloseNotify():
			h.hub.Unsubscribe(subscriber.ID)
			return
		case notification := <-subscriber.Ch:
			payloadBytes, err := json.Marshal(notification)

			if err != nil {
				log.Panicln("error marshaling payload")
			}

			c.Writer.WriteString("data: " + string(payloadBytes) + "\n\n")
			flusher.Flush()
		}
	}
}
