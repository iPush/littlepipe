package pipeline

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string
	Payload   any
	Metadata  map[string]any
	CreatedAt time.Time
	Error     error
}

func NewMessage(payload any) *Message {
	return &Message{
		ID:        uuid.New().String(),
		Payload:   payload,
		Metadata:  make(map[string]any),
		CreatedAt: time.Now(),
	}
}

func (m *Message) WithMetadata(key string, value any) *Message {
	m.Metadata[key] = value
	return m
}
