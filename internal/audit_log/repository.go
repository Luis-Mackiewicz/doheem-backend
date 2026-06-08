package audit_log

import (
	"context"

	"doheem-backend/internal/db"
)

type AuditLogRepo struct {
	q *db.Queries
}

func NewAuditLogRepo(q *db.Queries) *AuditLogRepo {
	return &AuditLogRepo{q: q}
}

func (r *AuditLogRepo) GetByID(ctx context.Context, id string) (AuditLog, error) {
	al, err := r.q.GetAuditLogByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return AuditLog{}, err
	}
	return toAuditLog(al), nil
}

func (r *AuditLogRepo) ListByGroup(ctx context.Context, groupID string) ([]AuditLogWithUser, error) {
	rows, err := r.q.ListAuditLogsByGroup(ctx, db.UUIDFromString(groupID))
	if err != nil {
		return nil, err
	}
	result := make([]AuditLogWithUser, len(rows))
	for i, row := range rows {
		result[i] = AuditLogWithUser{
			AuditLog: AuditLog{
				ID:         db.UUIDToString(row.ID),
				GroupID:    db.UUIDToStringPtr(row.GroupID),
				UserID:     db.UUIDToStringPtr(row.UserID),
				EntityType: row.EntityType,
				EntityID:   db.UUIDToString(row.EntityID),
				Action:     row.Action,
				Changes:    nil,
				CreatedAt:  row.CreatedAt.Time,
			},
			UserName: row.UserName.String,
		}
	}
	return result, nil
}

func (r *AuditLogRepo) ListByEntity(ctx context.Context, entityType, entityID string) ([]AuditLog, error) {
	logs, err := r.q.ListAuditLogsByEntity(ctx, db.ListAuditLogsByEntityParams{
		EntityType: entityType,
		EntityID:   db.UUIDFromString(entityID),
	})
	if err != nil {
		return nil, err
	}
	return toAuditLogs(logs), nil
}

func (r *AuditLogRepo) Create(ctx context.Context, groupID, userID, entityType, entityID, action string, changes map[string]interface{}) (AuditLog, error) {
	al, err := r.q.CreateAuditLog(ctx, db.CreateAuditLogParams{
		GroupID:    db.UUIDFromString(groupID),
		UserID:     db.UUIDFromString(userID),
		EntityType: entityType,
		EntityID:   db.UUIDFromString(entityID),
		Action:     action,
		Changes:    db.MapToJSON(changes),
	})
	if err != nil {
		return AuditLog{}, err
	}
	return toAuditLog(al), nil
}

func toAuditLog(al db.AuditLog) AuditLog {
	return AuditLog{
		ID:         db.UUIDToString(al.ID),
		GroupID:    db.UUIDToStringPtr(al.GroupID),
		UserID:     db.UUIDToStringPtr(al.UserID),
		EntityType: al.EntityType,
		EntityID:   db.UUIDToString(al.EntityID),
		Action:     al.Action,
		Changes:    db.JSONToMap(al.Changes),
		CreatedAt:  al.CreatedAt.Time,
	}
}

func toAuditLogs(logs []db.AuditLog) []AuditLog {
	result := make([]AuditLog, len(logs))
	for i, l := range logs {
		result[i] = toAuditLog(l)
	}
	return result
}
