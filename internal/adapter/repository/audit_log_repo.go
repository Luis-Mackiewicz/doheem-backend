package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type AuditLogRepo struct {
	q *db.Queries
}

func NewAuditLogRepo(q *db.Queries) *AuditLogRepo {
	return &AuditLogRepo{q: q}
}

func (r *AuditLogRepo) GetByID(ctx context.Context, id string) (domain.AuditLog, error) {
	al, err := r.q.GetAuditLogByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.AuditLog{}, err
	}
	return domainAuditLog(al), nil
}

func (r *AuditLogRepo) ListByGroup(ctx context.Context, groupID string) ([]domain.AuditLogWithUser, error) {
	rows, err := r.q.ListAuditLogsByGroup(ctx, uuidFromString(groupID))
	if err != nil {
		return nil, err
	}
	result := make([]domain.AuditLogWithUser, len(rows))
	for i, row := range rows {
		result[i] = domain.AuditLogWithUser{
			AuditLog: domain.AuditLog{
				ID:         uuidToString(row.ID),
				GroupID:    uuidToStringPtr(row.GroupID),
				UserID:     uuidToStringPtr(row.UserID),
				EntityType: row.EntityType,
				EntityID:   uuidToString(row.EntityID),
				Action:     row.Action,
				Changes:    nil,
				CreatedAt:  row.CreatedAt.Time,
			},
			UserName: row.UserName.String,
		}
	}
	return result, nil
}

func (r *AuditLogRepo) ListByEntity(ctx context.Context, entityType, entityID string) ([]domain.AuditLog, error) {
	logs, err := r.q.ListAuditLogsByEntity(ctx, db.ListAuditLogsByEntityParams{
		EntityType: entityType,
		EntityID:   uuidFromString(entityID),
	})
	if err != nil {
		return nil, err
	}
	return domainAuditLogs(logs), nil
}

func (r *AuditLogRepo) Create(ctx context.Context, groupID, userID, entityType, entityID, action string, changes map[string]interface{}) (domain.AuditLog, error) {
	al, err := r.q.CreateAuditLog(ctx, db.CreateAuditLogParams{
		GroupID:    uuidFromString(groupID),
		UserID:     uuidFromString(userID),
		EntityType: entityType,
		EntityID:   uuidFromString(entityID),
		Action:     action,
		Changes:    mapToJSON(changes),
	})
	if err != nil {
		return domain.AuditLog{}, err
	}
	return domainAuditLog(al), nil
}
