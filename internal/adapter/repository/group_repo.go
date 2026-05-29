package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type GroupRepo struct {
	q *db.Queries
}

func NewGroupRepo(q *db.Queries) *GroupRepo {
	return &GroupRepo{q: q}
}

func (r *GroupRepo) GetByID(ctx context.Context, id string) (domain.Group, error) {
	g, err := r.q.GetGroupByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.Group{}, err
	}
	return domainGroup(g), nil
}

func (r *GroupRepo) ListByUserID(ctx context.Context, userID string) ([]domain.Group, error) {
	groups, err := r.q.ListGroupsByUserID(ctx, uuidFromString(userID))
	if err != nil {
		return nil, err
	}
	return domainGroups(groups), nil
}

func (r *GroupRepo) Create(ctx context.Context, params domain.CreateGroupParams) (domain.Group, error) {
	g, err := r.q.CreateGroup(ctx, db.CreateGroupParams{
		Name:     params.Name,
		Currency: params.Currency,
	})
	if err != nil {
		return domain.Group{}, err
	}
	return domainGroup(g), nil
}

func (r *GroupRepo) Update(ctx context.Context, id string, params domain.UpdateGroupParams) (domain.Group, error) {
	var name string
	if params.Name != nil {
		name = *params.Name
	}
	var currency string
	if params.Currency != nil {
		currency = *params.Currency
	}
	g, err := r.q.UpdateGroup(ctx, db.UpdateGroupParams{
		ID:       uuidFromString(id),
		Name:     name,
		Currency: currency,
	})
	if err != nil {
		return domain.Group{}, err
	}
	return domainGroup(g), nil
}

func (r *GroupRepo) SoftDelete(ctx context.Context, id string) error {
	return r.q.SoftDeleteGroup(ctx, uuidFromString(id))
}

func (r *GroupRepo) Deactivate(ctx context.Context, id string) error {
	return r.q.DeactivateGroup(ctx, uuidFromString(id))
}

func (r *GroupRepo) Activate(ctx context.Context, id string) error {
	return r.q.ActivateGroup(ctx, uuidFromString(id))
}
