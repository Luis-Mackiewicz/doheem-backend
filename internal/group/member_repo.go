package group

import (
	"context"

	"doheem-backend/internal/db"
)

type GroupMemberRepo struct {
	q *db.Queries
}

func NewGroupMemberRepo(q *db.Queries) *GroupMemberRepo {
	return &GroupMemberRepo{q: q}
}

func (r *GroupMemberRepo) GetByID(ctx context.Context, id string) (GroupMember, error) {
	gm, err := r.q.GetGroupMemberByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return GroupMember{}, err
	}
	return toGroupMember(gm), nil
}

func (r *GroupMemberRepo) Get(ctx context.Context, groupID, userID string) (GroupMember, error) {
	gm, err := r.q.GetGroupMember(ctx, db.GetGroupMemberParams{
		GroupID: db.UUIDFromString(groupID),
		UserID:  db.UUIDFromString(userID),
	})
	if err != nil {
		return GroupMember{}, err
	}
	return toGroupMember(gm), nil
}

func (r *GroupMemberRepo) ListByGroup(ctx context.Context, groupID string) ([]GroupMemberWithUser, error) {
	rows, err := r.q.ListGroupMembers(ctx, db.UUIDFromString(groupID))
	if err != nil {
		return nil, err
	}
	result := make([]GroupMemberWithUser, len(rows))
	for i, row := range rows {
		result[i] = toGroupMemberWithUser(row)
	}
	return result, nil
}

func (r *GroupMemberRepo) Create(ctx context.Context, groupID, userID, role string) (GroupMember, error) {
	gm, err := r.q.CreateGroupMember(ctx, db.CreateGroupMemberParams{
		GroupID: db.UUIDFromString(groupID),
		UserID:  db.UUIDFromString(userID),
		Role:    role,
	})
	if err != nil {
		return GroupMember{}, err
	}
	return toGroupMember(gm), nil
}

func (r *GroupMemberRepo) UpdateRole(ctx context.Context, groupID, userID, role string) (GroupMember, error) {
	gm, err := r.q.UpdateGroupMemberRole(ctx, db.UpdateGroupMemberRoleParams{
		GroupID: db.UUIDFromString(groupID),
		UserID:  db.UUIDFromString(userID),
		Role:    role,
	})
	if err != nil {
		return GroupMember{}, err
	}
	return toGroupMember(gm), nil
}

func (r *GroupMemberRepo) Remove(ctx context.Context, groupID, userID string) error {
	return r.q.RemoveGroupMember(ctx, db.RemoveGroupMemberParams{
		GroupID: db.UUIDFromString(groupID),
		UserID:  db.UUIDFromString(userID),
	})
}

func (r *GroupMemberRepo) CountActive(ctx context.Context, groupID string) (int64, error) {
	return r.q.CountActiveGroupMembers(ctx, db.UUIDFromString(groupID))
}
