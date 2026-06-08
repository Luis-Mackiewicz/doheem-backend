package group

import (
	"context"

	"doheem-backend/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
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
	g, err := r.q.CreateGroup(ctx, params.Name)
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
		Cnpj:        deptrStr(params.Cnpj),
		Cep:         deptrStr(params.Cep),
		PhotoUrl:    db.TextFromStringPtr(params.PhotoURL),
	})
	if err != nil {
		return Group{}, err
	}
	return toGroup(g), nil
}

func (r *GroupRepo) RegenerateInviteToken(ctx context.Context, id, token string) error {
	_, err := r.q.RegenerateInviteToken(ctx, db.RegenerateInviteTokenParams{
		ID:          db.UUIDFromString(id),
		InviteToken: pgtype.Text{String: token, Valid: true},
	})
	return err
}

func toGroup(g db.Group) Group {
	return Group{
		ID:          db.UUIDToString(g.ID),
		Name:        g.Name,
		Description: g.Description,
		MonthlyFee:  db.NumericToFloat64(g.MonthlyFee),
		Cnpj:        g.Cnpj,
		Cep:         g.Cep,
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

func deptrNumeric(f *float64) pgtype.Numeric {
	if f != nil {
		return db.NumericFromFloat64(*f)
	}
	return pgtype.Numeric{}
}
