package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type GroupMemberRepo struct {
	q *db.Queries
}

func NewGroupMemberRepo(q *db.Queries) *GroupMemberRepo {
	return &GroupMemberRepo{q: q}
}

func (r *GroupMemberRepo) GetByID(ctx context.Context, id string) (domain.GroupMember, error) {
	gm, err := r.q.GetGroupMemberByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.GroupMember{}, err
	}
	return domainGroupMember(gm), nil
}

func (r *GroupMemberRepo) Get(ctx context.Context, groupID, userID string) (domain.GroupMember, error) {
	gm, err := r.q.GetGroupMember(ctx, db.GetGroupMemberParams{
		GroupID: uuidFromString(groupID),
		UserID:  uuidFromString(userID),
	})
	if err != nil {
		return domain.GroupMember{}, err
	}
	return domainGroupMember(gm), nil
}

func (r *GroupMemberRepo) ListByGroup(ctx context.Context, groupID string) ([]domain.GroupMemberWithUser, error) {
	rows, err := r.q.ListGroupMembers(ctx, uuidFromString(groupID))
	if err != nil {
		return nil, err
	}
	result := make([]domain.GroupMemberWithUser, len(rows))
	for i, row := range rows {
		result[i] = domainGroupMemberWithUser(row)
	}
	return result, nil
}

func (r *GroupMemberRepo) Create(ctx context.Context, groupID, userID, role string) (domain.GroupMember, error) {
	gm, err := r.q.CreateGroupMember(ctx, db.CreateGroupMemberParams{
		GroupID: uuidFromString(groupID),
		UserID:  uuidFromString(userID),
		Role:    role,
	})
	if err != nil {
		return domain.GroupMember{}, err
	}
	return domainGroupMember(gm), nil
}

func (r *GroupMemberRepo) UpdateRole(ctx context.Context, groupID, userID, role string) (domain.GroupMember, error) {
	gm, err := r.q.UpdateGroupMemberRole(ctx, db.UpdateGroupMemberRoleParams{
		GroupID: uuidFromString(groupID),
		UserID:  uuidFromString(userID),
		Role:    role,
	})
	if err != nil {
		return domain.GroupMember{}, err
	}
	return domainGroupMember(gm), nil
}

func (r *GroupMemberRepo) Remove(ctx context.Context, groupID, userID string) error {
	return r.q.RemoveGroupMember(ctx, db.RemoveGroupMemberParams{
		GroupID: uuidFromString(groupID),
		UserID:  uuidFromString(userID),
	})
}

func (r *GroupMemberRepo) CountActive(ctx context.Context, groupID string) (int64, error) {
	return r.q.CountActiveGroupMembers(ctx, uuidFromString(groupID))
}
