package audit_log

import (
	"context"
	"time"
)

type AuditLog struct {
	ID         string
	GroupID    *string
	UserID     *string
	EntityType string
	EntityID   string
	Action     string
	Changes    map[string]interface{}
	CreatedAt  time.Time
}

type AuditLogWithUser struct {
	AuditLog
	UserName string
}

type AuditLogRepository interface {
	GetByID(ctx context.Context, id string) (AuditLog, error)
	ListByGroup(ctx context.Context, groupID string) ([]AuditLogWithUser, error)
	ListByEntity(ctx context.Context, entityType, entityID string) ([]AuditLog, error)
	Create(ctx context.Context, groupID, userID, entityType, entityID, action string, changes map[string]interface{}) (AuditLog, error)
}
