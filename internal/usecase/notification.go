package usecase

import (
	"github.com/google/uuid"
)

type Notification struct {
	Type    string          `json:"type"`
	Payload interface{}     `json:"payload"`
	UserID  string          `json:"userID,omitempty"`
}

type Subscriber struct {
	ID     string
	UserID string
	Ch     chan Notification
}

type NotificationHub struct {
	subscribers map[string]*Subscriber
	lock        chan struct{}
}

func NewNotificationHub() *NotificationHub {
	return &NotificationHub{
		subscribers: make(map[string]*Subscriber),
		lock:        make(chan struct{}, 1),
	}
}

func (h *NotificationHub) Subscribe(userID string) (*Subscriber, error) {
	h.lock <- struct{}{}
	defer func() { <-h.lock }()

	id := uuid.New().String()
	subscriber := &Subscriber{
		ID:     id,
		UserID: userID,
		Ch:     make(chan Notification),
	}

	h.subscribers[id] = subscriber
	return subscriber, nil
}

func (h *NotificationHub) Unsubscribe(subscriberID string) {
	h.lock <- struct{}{}
	defer func() { <-h.lock }()

	if sub, exists := h.subscribers[subscriberID]; exists {
		close(sub.Ch)
		delete(h.subscribers, subscriberID)
	}
}

func (h *NotificationHub) Notify(notification Notification) {
	h.lock <- struct{}{}
	defer func() { <-h.lock }()

	for _, sub := range h.subscribers {
		if notification.UserID == "" || sub.UserID == notification.UserID {
			select {
			case sub.Ch <- notification:
			default:
				// If channel is full, drop the notification
			}
		}
	}
}
