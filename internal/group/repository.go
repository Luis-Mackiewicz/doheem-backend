package group

import (
	"context"

	"doheem-backend/internal/db"
)

type GroupRepo struct {
	q *db.Queries
}

func NewGroupRepo(q *db.Queries) *GroupRepo {
	return &GroupRepo{q: q}
}

func (r *GroupRepo) GetByID(ctx context.Context, id string) (Group, error) {
	g, err := r.q.GetGroupByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return Group{}, err
	}
	return toGroup(g), nil
}

func (r *GroupRepo) ListByUserID(ctx context.Context, userID string) ([]Group, error) {
	groups, err := r.q.ListGroupsByUserID(ctx, db.UUIDFromString(userID))
	if err != nil {
		return nil, err
	}
	return toGroups(groups), nil
}

func (r *GroupRepo) Create(ctx context.Context, params CreateGroupParams) (Group, error) {
	g, err := r.q.CreateGroup(ctx, db.CreateGroupParams{
		Name:     params.Name,
		Currency: params.Currency,
	})
	if err != nil {
		return Group{}, err
	}
	return toGroup(g), nil
}

func (r *GroupRepo) Update(ctx context.Context, id string, params UpdateGroupParams) (Group, error) {
	var name string
	if params.Name != nil {
		name = *params.Name
	}
	var currency string
	if params.Currency != nil {
		currency = *params.Currency
	}
	g, err := r.q.UpdateGroup(ctx, db.UpdateGroupParams{
		ID:       db.UUIDFromString(id),
		Name:     name,
		Currency: currency,
	})
	if err != nil {
		return Group{}, err
	}
	return toGroup(g), nil
}

func (r *GroupRepo) SoftDelete(ctx context.Context, id string) error {
	return r.q.SoftDeleteGroup(ctx, db.UUIDFromString(id))
}

func (r *GroupRepo) Deactivate(ctx context.Context, id string) error {
	return r.q.DeactivateGroup(ctx, db.UUIDFromString(id))
}

func (r *GroupRepo) Activate(ctx context.Context, id string) error {
	return r.q.ActivateGroup(ctx, db.UUIDFromString(id))
}

func toGroup(g db.Group) Group {
	return Group{
		ID:            db.UUIDToString(g.ID),
		Name:          g.Name,
		Currency:      g.Currency,
		IsActive:      g.IsActive,
		InactiveSince: db.TimestamptzToTimePtr(g.InactiveSince),
		CreatedAt:     g.CreatedAt.Time,
		UpdatedAt:     g.UpdatedAt.Time,
		DeletedAt:     db.TimestamptzToTimePtr(g.DeletedAt),
	}
}

func toGroups(groups []db.Group) []Group {
	result := make([]Group, len(groups))
	for i, g := range groups {
		result[i] = toGroup(g)
	}
	return result
}

func toGroupMember(gm db.GroupMember) GroupMember {
	return GroupMember{
		ID:       db.UUIDToString(gm.ID),
		GroupID:  db.UUIDToString(gm.GroupID),
		UserID:   db.UUIDToString(gm.UserID),
		Role:     gm.Role,
		JoinedAt: gm.JoinedAt.Time,
		LeftAt:   db.TimestamptzToTimePtr(gm.LeftAt),
		IsActive: gm.IsActive,
	}
}

func toGroupMemberWithUser(row db.ListGroupMembersRow) GroupMemberWithUser {
	return GroupMemberWithUser{
		GroupMember: GroupMember{
			ID:       db.UUIDToString(row.ID),
			GroupID:  db.UUIDToString(row.GroupID),
			UserID:   db.UUIDToString(row.UserID),
			Role:     row.Role,
			JoinedAt: row.JoinedAt.Time,
			LeftAt:   db.TimestamptzToTimePtr(row.LeftAt),
			IsActive: row.IsActive,
		},
		UserName:  row.Name,
		UserEmail: row.Email,
		AvatarURL: db.TextToStringPtr(row.AvatarUrl),
	}
}
