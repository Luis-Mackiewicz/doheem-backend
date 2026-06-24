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

func (r *GroupMemberRepo) Create(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error) {
	gm, err := r.q.CreateGroupMember(ctx, db.CreateGroupMemberParams{
		GroupID: db.UUIDFromString(groupID),
		UserID:  db.UUIDFromString(userID),
		IsAdmin: isAdmin,
	})
	if err != nil {
		return GroupMember{}, err
	}
	return toGroupMember(gm), nil
}

func (r *GroupMemberRepo) UpdateRole(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error) {
	gm, err := r.q.UpdateGroupMemberRole(ctx, db.UpdateGroupMemberRoleParams{
		GroupID: db.UUIDFromString(groupID),
		UserID:  db.UUIDFromString(userID),
		IsAdmin: isAdmin,
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

func (r *GroupMemberRepo) Count(ctx context.Context, groupID string) (int64, error) {
	return r.q.CountGroupMembers(ctx, db.UUIDFromString(groupID))
}

func (r *GroupMemberRepo) CountAdmins(ctx context.Context, groupID string) (int64, error) {
	return r.q.CountGroupAdmins(ctx, db.UUIDFromString(groupID))
}

func toGroupMember(gm db.GroupMember) GroupMember {
	return GroupMember{
		ID:       db.UUIDToString(gm.ID),
		GroupID:  db.UUIDToString(gm.GroupID),
		UserID:   db.UUIDToString(gm.UserID),
		IsAdmin:  gm.IsAdmin,
		JoinedAt: gm.JoinedAt.Time,
	}
}

func toGroupMemberWithUser(row db.ListGroupMembersRow) GroupMemberWithUser {
	return GroupMemberWithUser{
		GroupMember: GroupMember{
			ID:       db.UUIDToString(row.ID),
			GroupID:  db.UUIDToString(row.GroupID),
			UserID:   db.UUIDToString(row.UserID),
			IsAdmin:  row.IsAdmin,
			JoinedAt: row.JoinedAt.Time,
		},
		UserName:  row.Name,
		UserEmail: row.Email,
		UserPhone: db.TextToStringPtr(row.Phone),
		AvatarURL: db.TextToStringPtr(row.AvatarUrl),
	}
}
