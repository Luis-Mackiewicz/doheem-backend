package expense

import (
	"context"

	"doheem-backend/internal/db"
)

type ExpenseCategoryRepo struct {
	q *db.Queries
}

func NewExpenseCategoryRepo(q *db.Queries) *ExpenseCategoryRepo {
	return &ExpenseCategoryRepo{q: q}
}

func (r *ExpenseCategoryRepo) GetByID(ctx context.Context, id string) (ExpenseCategory, error) {
	ec, err := r.q.GetExpenseCategoryByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return ExpenseCategory{}, err
	}
	return toCategory(ec), nil
}

func (r *ExpenseCategoryRepo) ListByGroup(ctx context.Context, groupID string) ([]ExpenseCategory, error) {
	categories, err := r.q.ListExpenseCategoriesByGroup(ctx, db.UUIDFromString(groupID))
	if err != nil {
		return nil, err
	}
	return toCategories(categories), nil
}

func (r *ExpenseCategoryRepo) Create(ctx context.Context, groupID, name string) (ExpenseCategory, error) {
	ec, err := r.q.CreateExpenseCategory(ctx, db.CreateExpenseCategoryParams{
		GroupID: db.UUIDFromString(groupID),
		Name:    name,
	})
	if err != nil {
		return ExpenseCategory{}, err
	}
	return toCategory(ec), nil
}

func (r *ExpenseCategoryRepo) Update(ctx context.Context, id, groupID, name string) (ExpenseCategory, error) {
	ec, err := r.q.UpdateExpenseCategory(ctx, db.UpdateExpenseCategoryParams{
		ID:      db.UUIDFromString(id),
		GroupID: db.UUIDFromString(groupID),
		Name:    name,
	})
	if err != nil {
		return ExpenseCategory{}, err
	}
	return toCategory(ec), nil
}

func (r *ExpenseCategoryRepo) Delete(ctx context.Context, id, groupID string) error {
	return r.q.DeleteExpenseCategory(ctx, db.DeleteExpenseCategoryParams{
		ID:      db.UUIDFromString(id),
		GroupID: db.UUIDFromString(groupID),
	})
}

func toCategory(ec db.ExpenseCategory) ExpenseCategory {
	return ExpenseCategory{
		ID:        db.UUIDToString(ec.ID),
		GroupID:   db.UUIDToString(ec.GroupID),
		Name:      ec.Name,
		CreatedAt: ec.CreatedAt.Time,
	}
}

func toCategories(categories []db.ExpenseCategory) []ExpenseCategory {
	result := make([]ExpenseCategory, len(categories))
	for i, c := range categories {
		result[i] = toCategory(c)
	}
	return result
}
