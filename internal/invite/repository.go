package invite

import (
	"context"
	"time"

	"doheem-backend/internal/db"
)

type InviteRepo struct {
	q *db.Queries
}

func NewInviteRepo(q *db.Queries) *InviteRepo {
	return &InviteRepo{q: q}
}

func (r *InviteRepo) GetByID(ctx context.Context, id string) (Invite, error) {
	i, err := r.q.GetInviteByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return Invite{}, err
	}
	return toInvite(i), nil
}

func (r *InviteRepo) GetByCode(ctx context.Context, code string) (InviteWithGroup, error) {
	row, err := r.q.GetInviteByCode(ctx, code)
	if err != nil {
		return InviteWithGroup{}, err
	}
	return toInviteWithGroup(row), nil
}

func (r *InviteRepo) ListByGroup(ctx context.Context, groupID string) ([]InviteWithCreator, error) {
	rows, err := r.q.ListInvitesByGroup(ctx, db.UUIDFromString(groupID))
	if err != nil {
		return nil, err
	}
	result := make([]InviteWithCreator, len(rows))
	for i, row := range rows {
		result[i] = InviteWithCreator{
			Invite: Invite{
				ID:        db.UUIDToString(row.ID),
				GroupID:   db.UUIDToString(row.GroupID),
				Code:      row.Code,
				CreatedBy: db.UUIDToString(row.CreatedBy),
				ExpiresAt: row.ExpiresAt.Time,
				UsedAt:    db.TimestamptzToTimePtr(row.UsedAt),
				RevokedAt: db.TimestamptzToTimePtr(row.RevokedAt),
				CreatedAt: row.CreatedAt.Time,
			},
			CreatedByName: row.CreatedByName,
		}
	}
	return result, nil
}

func (r *InviteRepo) ListPendingByUser(ctx context.Context, userID string) ([]InviteWithGroup, error) {
	rows, err := r.q.ListPendingInvitesByUser(ctx, db.UUIDFromString(userID))
	if err != nil {
		return nil, err
	}
	result := make([]InviteWithGroup, len(rows))
	for i, row := range rows {
		result[i] = InviteWithGroup{
			Invite: Invite{
				ID:        db.UUIDToString(row.ID),
				GroupID:   db.UUIDToString(row.GroupID),
				Code:      row.Code,
				CreatedBy: db.UUIDToString(row.CreatedBy),
				ExpiresAt: row.ExpiresAt.Time,
				UsedAt:    db.TimestamptzToTimePtr(row.UsedAt),
				RevokedAt: db.TimestamptzToTimePtr(row.RevokedAt),
				CreatedAt: row.CreatedAt.Time,
			},
			GroupName: row.GroupName,
		}
	}
	return result, nil
}

func (r *InviteRepo) Create(ctx context.Context, groupID, code, createdBy string, expiresAt time.Time) (Invite, error) {
	i, err := r.q.CreateInvite(ctx, db.CreateInviteParams{
		GroupID:   db.UUIDFromString(groupID),
		Code:      code,
		CreatedBy: db.UUIDFromString(createdBy),
		ExpiresAt: db.TimestamptzFromTime(expiresAt),
	})
	if err != nil {
		return Invite{}, err
	}
	return toInvite(i), nil
}

func (r *InviteRepo) Use(ctx context.Context, id string) error {
	return r.q.UseInvite(ctx, db.UUIDFromString(id))
}

func (r *InviteRepo) Revoke(ctx context.Context, id string) error {
	return r.q.RevokeInvite(ctx, db.UUIDFromString(id))
}

func toInvite(i db.Invite) Invite {
	return Invite{
		ID:        db.UUIDToString(i.ID),
		GroupID:   db.UUIDToString(i.GroupID),
		Code:      i.Code,
		CreatedBy: db.UUIDToString(i.CreatedBy),
		ExpiresAt: i.ExpiresAt.Time,
		UsedAt:    db.TimestamptzToTimePtr(i.UsedAt),
		RevokedAt: db.TimestamptzToTimePtr(i.RevokedAt),
		CreatedAt: i.CreatedAt.Time,
	}
}

func toInviteWithGroup(row db.GetInviteByCodeRow) InviteWithGroup {
	return InviteWithGroup{
		Invite: Invite{
			ID:        db.UUIDToString(row.ID),
			GroupID:   db.UUIDToString(row.GroupID),
			Code:      row.Code,
			CreatedBy: db.UUIDToString(row.CreatedBy),
			ExpiresAt: row.ExpiresAt.Time,
			UsedAt:    db.TimestamptzToTimePtr(row.UsedAt),
			RevokedAt: db.TimestamptzToTimePtr(row.RevokedAt),
			CreatedAt: row.CreatedAt.Time,
		},
		GroupName: row.GroupName,
	}
}
