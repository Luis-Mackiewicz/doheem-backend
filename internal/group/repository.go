package group

import (
	"context"

	"doheem-backend/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
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
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		return Group{}, err
	}
	return toGroup(g), nil
}

func (r *GroupRepo) Update(ctx context.Context, id string, params UpdateGroupParams) (Group, error) {
	g, err := r.q.UpdateGroup(ctx, db.UpdateGroupParams{
		ID:          db.UUIDFromString(id),
		Name:        deptrStr(params.Name),
		Description: deptrStr(params.Description),
		MonthlyFee:  deptrNumeric(params.MonthlyFee),
		PhotoUrl:    db.TextFromStringPtr(params.PhotoURL),
	})
	if err != nil {
		return Group{}, err
	}
	return toGroup(g), nil
}

func (r *GroupRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeleteGroup(ctx, db.UUIDFromString(id))
}

func (r *GroupRepo) RegenerateInviteToken(ctx context.Context, id, token string) error {
	_, err := r.q.RegenerateInviteToken(ctx, db.RegenerateInviteTokenParams{
		ID:          db.UUIDFromString(id),
		InviteToken: pgtype.Text{String: token, Valid: true},
	})
	return err
}

func (r *GroupRepo) CountByUserID(ctx context.Context, userID string) (int64, error) {
	return r.q.CountGroupsByUserID(ctx, db.UUIDFromString(userID))
}

func toGroup(g db.Group) Group {
	return Group{
		ID:          db.UUIDToString(g.ID),
		Name:        g.Name,
		Description: g.Description,
		MonthlyFee:  db.NumericToDecimal(g.MonthlyFee),
		PhotoURL:    db.TextToStringPtr(g.PhotoUrl),
		InviteToken: db.TextToStringPtr(g.InviteToken),
		CreatedAt:   g.CreatedAt.Time,
		UpdatedAt:   g.UpdatedAt.Time,
	}
}

func toGroups(groups []db.Group) []Group {
	result := make([]Group, len(groups))
	for i, g := range groups {
		result[i] = toGroup(g)
	}
	return result
}

func deptrStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func deptrNumeric(d *decimal.Decimal) pgtype.Numeric {
	if d != nil {
		return db.NumericFromDecimal(*d)
	}
	return pgtype.Numeric{}
}
