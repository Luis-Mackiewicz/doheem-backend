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

func (r *ExpenseCategoryRepo) ListAll(ctx context.Context) ([]ExpenseCategory, error) {
	categories, err := r.q.ListExpenseCategories(ctx)
	if err != nil {
		return nil, err
	}
	return toCategories(categories), nil
}

func (r *ExpenseCategoryRepo) Create(ctx context.Context, slug, label string) (ExpenseCategory, error) {
	ec, err := r.q.CreateExpenseCategory(ctx, db.CreateExpenseCategoryParams{
		Slug:  slug,
		Label: label,
	})
	if err != nil {
		return ExpenseCategory{}, err
	}
	return toCategory(ec), nil
}

func (r *ExpenseCategoryRepo) Update(ctx context.Context, id, slug, label string) (ExpenseCategory, error) {
	ec, err := r.q.UpdateExpenseCategory(ctx, db.UpdateExpenseCategoryParams{
		ID:    db.UUIDFromString(id),
		Slug:  slug,
		Label: label,
	})
	if err != nil {
		return ExpenseCategory{}, err
	}
	return toCategory(ec), nil
}

func (r *ExpenseCategoryRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeleteExpenseCategory(ctx, db.UUIDFromString(id))
}

func toCategory(ec db.ExpenseCategory) ExpenseCategory {
	return ExpenseCategory{
		ID:        db.UUIDToString(ec.ID),
		Slug:      ec.Slug,
		Label:     ec.Label,
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
