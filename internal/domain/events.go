package domain

import (
	"context"
	"time"
)

type DomainEvent struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	EntityID  string         `json:"entity_id"`
	UserID    string         `json:"user_id,omitempty"`
	GroupID   string         `json:"group_id,omitempty"`
	Payload   map[string]any `json:"payload,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

type EventBus interface {
	Publish(ctx context.Context, event DomainEvent) error
}
