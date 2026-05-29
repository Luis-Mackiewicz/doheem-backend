package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
)

type ExpenseRepo struct {
	q *db.Queries
}

func NewExpenseRepo(q *db.Queries) *ExpenseRepo {
	return &ExpenseRepo{q: q}
}

func (r *ExpenseRepo) GetByID(ctx context.Context, id string) (domain.Expense, error) {
	e, err := r.q.GetExpenseByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.Expense{}, err
	}
	return domainExpense(e), nil
}

func (r *ExpenseRepo) ListByGroup(ctx context.Context, groupID string) ([]domain.Expense, error) {
	expenses, err := r.q.ListExpensesByGroup(ctx, uuidFromString(groupID))
	if err != nil {
		return nil, err
	}
	return domainExpenses(expenses), nil
}

func (r *ExpenseRepo) ListByUser(ctx context.Context, userID string) ([]domain.Expense, error) {
	expenses, err := r.q.ListExpensesByUser(ctx, uuidFromString(userID))
	if err != nil {
		return nil, err
	}
	return domainExpenses(expenses), nil
}

func (r *ExpenseRepo) ListByCategory(ctx context.Context, categoryID string) ([]domain.Expense, error) {
	expenses, err := r.q.ListExpensesByCategory(ctx, uuidFromString(categoryID))
	if err != nil {
		return nil, err
	}
	return domainExpenses(expenses), nil
}

func (r *ExpenseRepo) Create(ctx context.Context, params domain.CreateExpenseParams) (domain.Expense, error) {
	var dueDate pgtype.Date
	if params.DueDate != nil {
		dueDate = dateFromTime(*params.DueDate)
	}
	e, err := r.q.CreateExpense(ctx, db.CreateExpenseParams{
		GroupID:          uuidFromString(params.GroupID),
		CreatedBy:        uuidFromString(params.CreatedBy),
		Description:      params.Description,
		TotalAmount:      numericFromFloat64(params.TotalAmount),
		ExpenseDate:      dateFromTime(params.ExpenseDate),
		DueDate:          dueDate,
		CategoryID:       uuidFromStringPtr(params.CategoryID),
		SplitType:        params.SplitType,
		IsInstallment:    params.IsInstallment,
		InstallmentCount: int2FromInt16Ptr(params.InstallmentCount),
	})
	if err != nil {
		return domain.Expense{}, err
	}
	return domainExpense(e), nil
}

func (r *ExpenseRepo) Update(ctx context.Context, id string, params domain.UpdateExpenseParams) (domain.Expense, error) {
	var desc string
	if params.Description != nil {
		desc = *params.Description
	}
	var total pgtype.Numeric
	if params.TotalAmount != nil {
		total = numericFromFloat64(*params.TotalAmount)
	}
	var expDate pgtype.Date
	if params.ExpenseDate != nil {
		expDate = dateFromTime(*params.ExpenseDate)
	}
	var dueDate pgtype.Date
	if params.DueDate != nil {
		dueDate = dateFromTime(*params.DueDate)
	}
	var split string
	if params.SplitType != nil {
		split = *params.SplitType
	}
	e, err := r.q.UpdateExpense(ctx, db.UpdateExpenseParams{
		ID:          uuidFromString(id),
		Description: desc,
		TotalAmount: total,
		ExpenseDate: expDate,
		DueDate:     dueDate,
		CategoryID:  uuidFromStringPtrOrZero(params.CategoryID),
		SplitType:   split,
	})
	if err != nil {
		return domain.Expense{}, err
	}
	return domainExpense(e), nil
}

func (r *ExpenseRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeleteExpense(ctx, uuidFromString(id))
}

func (r *ExpenseRepo) GetTotalByGroup(ctx context.Context, groupID string) (float64, error) {
	val, err := r.q.GetTotalExpensesByGroup(ctx, uuidFromString(groupID))
	if err != nil {
		return 0, err
	}
	return numericToFloat64(val), nil
}
