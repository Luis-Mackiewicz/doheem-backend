package repository

import (
	"context"
	"time"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type InviteRepo struct {
	q *db.Queries
}

func NewInviteRepo(q *db.Queries) *InviteRepo {
	return &InviteRepo{q: q}
}

func (r *InviteRepo) GetByID(ctx context.Context, id string) (domain.Invite, error) {
	i, err := r.q.GetInviteByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.Invite{}, err
	}
	return domainInvite(i), nil
}

func (r *InviteRepo) GetByCode(ctx context.Context, code string) (domain.InviteWithGroup, error) {
	row, err := r.q.GetInviteByCode(ctx, code)
	if err != nil {
		return domain.InviteWithGroup{}, err
	}
	return domainInviteWithGroup(row), nil
}

func (r *InviteRepo) ListByGroup(ctx context.Context, groupID string) ([]domain.InviteWithCreator, error) {
	rows, err := r.q.ListInvitesByGroup(ctx, uuidFromString(groupID))
	if err != nil {
		return nil, err
	}
	result := make([]domain.InviteWithCreator, len(rows))
	for i, row := range rows {
		result[i] = domain.InviteWithCreator{
			Invite: domain.Invite{
				ID:        uuidToString(row.ID),
				GroupID:   uuidToString(row.GroupID),
				Code:      row.Code,
				CreatedBy: uuidToString(row.CreatedBy),
				ExpiresAt: row.ExpiresAt.Time,
				UsedAt:    timestamptzToTimePtr(row.UsedAt),
				RevokedAt: timestamptzToTimePtr(row.RevokedAt),
				CreatedAt: row.CreatedAt.Time,
			},
			CreatedByName: row.CreatedByName,
		}
	}
	return result, nil
}

func (r *InviteRepo) ListPendingByUser(ctx context.Context, userID string) ([]domain.InviteWithGroup, error) {
	rows, err := r.q.ListPendingInvitesByUser(ctx, uuidFromString(userID))
	if err != nil {
		return nil, err
	}
	result := make([]domain.InviteWithGroup, len(rows))
	for i, row := range rows {
		result[i] = domain.InviteWithGroup{
			Invite: domain.Invite{
				ID:        uuidToString(row.ID),
				GroupID:   uuidToString(row.GroupID),
				Code:      row.Code,
				CreatedBy: uuidToString(row.CreatedBy),
				ExpiresAt: row.ExpiresAt.Time,
				UsedAt:    timestamptzToTimePtr(row.UsedAt),
				RevokedAt: timestamptzToTimePtr(row.RevokedAt),
				CreatedAt: row.CreatedAt.Time,
			},
			GroupName: row.GroupName,
		}
	}
	return result, nil
}

func (r *InviteRepo) Create(ctx context.Context, groupID, code, createdBy string, expiresAt time.Time) (domain.Invite, error) {
	i, err := r.q.CreateInvite(ctx, db.CreateInviteParams{
		GroupID:   uuidFromString(groupID),
		Code:      code,
		CreatedBy: uuidFromString(createdBy),
		ExpiresAt: timestamptzFromTime(expiresAt),
	})
	if err != nil {
		return domain.Invite{}, err
	}
	return domainInvite(i), nil
}

func (r *InviteRepo) Use(ctx context.Context, id string) error {
	return r.q.UseInvite(ctx, uuidFromString(id))
}

func (r *InviteRepo) Revoke(ctx context.Context, id string) error {
	return r.q.RevokeInvite(ctx, uuidFromString(id))
}
