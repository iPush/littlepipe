package pipeline

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string
	Payload   *Record
	Metadata  map[string]any
	CreatedAt time.Time
	Error     error

	TraceID string
	SpanID  string
	Context context.Context
}

func NewMessage(payload *Record) *Message {
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
