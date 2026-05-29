package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type ExpenseCategoryRepo struct {
	q *db.Queries
}

func NewExpenseCategoryRepo(q *db.Queries) *ExpenseCategoryRepo {
	return &ExpenseCategoryRepo{q: q}
}

func (r *ExpenseCategoryRepo) GetByID(ctx context.Context, id string) (domain.ExpenseCategory, error) {
	ec, err := r.q.GetExpenseCategoryByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.ExpenseCategory{}, err
	}
	return domainExpenseCategory(ec), nil
}

func (r *ExpenseCategoryRepo) ListByGroup(ctx context.Context, groupID string) ([]domain.ExpenseCategory, error) {
	categories, err := r.q.ListExpenseCategoriesByGroup(ctx, uuidFromString(groupID))
	if err != nil {
		return nil, err
	}
	return domainExpenseCategories(categories), nil
}

func (r *ExpenseCategoryRepo) Create(ctx context.Context, groupID, name string) (domain.ExpenseCategory, error) {
	ec, err := r.q.CreateExpenseCategory(ctx, db.CreateExpenseCategoryParams{
		GroupID: uuidFromString(groupID),
		Name:    name,
	})
	if err != nil {
		return domain.ExpenseCategory{}, err
	}
	return domainExpenseCategory(ec), nil
}

func (r *ExpenseCategoryRepo) Update(ctx context.Context, id, groupID, name string) (domain.ExpenseCategory, error) {
	ec, err := r.q.UpdateExpenseCategory(ctx, db.UpdateExpenseCategoryParams{
		ID:      uuidFromString(id),
		GroupID: uuidFromString(groupID),
		Name:    name,
	})
	if err != nil {
		return domain.ExpenseCategory{}, err
	}
	return domainExpenseCategory(ec), nil
}

func (r *ExpenseCategoryRepo) Delete(ctx context.Context, id, groupID string) error {
	return r.q.DeleteExpenseCategory(ctx, db.DeleteExpenseCategoryParams{
		ID:      uuidFromString(id),
		GroupID: uuidFromString(groupID),
	})
}
