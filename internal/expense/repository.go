package expense

import (
	"context"

	"doheem-backend/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type ExpenseRepo struct {
	q *db.Queries
}

func NewExpenseRepo(q *db.Queries) *ExpenseRepo {
	return &ExpenseRepo{q: q}
}

func (r *ExpenseRepo) GetByID(ctx context.Context, id string) (Expense, error) {
	e, err := r.q.GetExpenseByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return Expense{}, err
	}
	return toExpense(e), nil
}

func (r *ExpenseRepo) ListByGroup(ctx context.Context, groupID string) ([]Expense, error) {
	expenses, err := r.q.ListExpensesByGroup(ctx, db.UUIDFromString(groupID))
	if err != nil {
		return nil, err
	}
	return toExpenses(expenses), nil
}

func (r *ExpenseRepo) ListByUser(ctx context.Context, userID string) ([]Expense, error) {
	expenses, err := r.q.ListExpensesByUser(ctx, db.UUIDFromString(userID))
	if err != nil {
		return nil, err
	}
	return toExpenses(expenses), nil
}

func (r *ExpenseRepo) ListByCategory(ctx context.Context, categoryID string) ([]Expense, error) {
	expenses, err := r.q.ListExpensesByCategory(ctx, db.UUIDFromString(categoryID))
	if err != nil {
		return nil, err
	}
	return toExpenses(expenses), nil
}

func (r *ExpenseRepo) Create(ctx context.Context, params CreateExpenseParams) (Expense, error) {
	var dueDate pgtype.Date
	if params.DueDate != nil {
		dueDate = db.DateFromTime(*params.DueDate)
	}
	e, err := r.q.CreateExpense(ctx, db.CreateExpenseParams{
		GroupID:          db.UUIDFromString(params.GroupID),
		CreatedBy:        db.UUIDFromString(params.CreatedBy),
		Description:      params.Description,
		TotalAmount:      db.NumericFromFloat64(params.TotalAmount),
		ExpenseDate:      db.DateFromTime(params.ExpenseDate),
		DueDate:          dueDate,
		CategoryID:       db.UUIDFromStringPtr(params.CategoryID),
		SplitType:        params.SplitType,
		IsInstallment:    params.IsInstallment,
		InstallmentCount: db.Int2FromInt16Ptr(params.InstallmentCount),
	})
	if err != nil {
		return Expense{}, err
	}
	return toExpense(e), nil
}

func (r *ExpenseRepo) Update(ctx context.Context, id string, params UpdateExpenseParams) (Expense, error) {
	var desc string
	if params.Description != nil {
		desc = *params.Description
	}
	var total pgtype.Numeric
	if params.TotalAmount != nil {
		total = db.NumericFromFloat64(*params.TotalAmount)
	}
	var expDate pgtype.Date
	if params.ExpenseDate != nil {
		expDate = db.DateFromTime(*params.ExpenseDate)
	}
	var dueDate pgtype.Date
	if params.DueDate != nil {
		dueDate = db.DateFromTime(*params.DueDate)
	}
	var split string
	if params.SplitType != nil {
		split = *params.SplitType
	}
	e, err := r.q.UpdateExpense(ctx, db.UpdateExpenseParams{
		ID:          db.UUIDFromString(id),
		Description: desc,
		TotalAmount: total,
		ExpenseDate: expDate,
		DueDate:     dueDate,
		CategoryID:  db.UUIDFromStringPtr(params.CategoryID),
		SplitType:   split,
	})
	if err != nil {
		return Expense{}, err
	}
	return toExpense(e), nil
}

func (r *ExpenseRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeleteExpense(ctx, db.UUIDFromString(id))
}

func (r *ExpenseRepo) GetTotalByGroup(ctx context.Context, groupID string) (float64, error) {
	val, err := r.q.GetTotalExpensesByGroup(ctx, db.UUIDFromString(groupID))
	if err != nil {
		return 0, err
	}
	return db.NumericToFloat64(val), nil
}

func toExpense(e db.Expense) Expense {
	return Expense{
		ID:               db.UUIDToString(e.ID),
		GroupID:          db.UUIDToString(e.GroupID),
		CreatedBy:        db.UUIDToString(e.CreatedBy),
		Description:      e.Description,
		TotalAmount:      db.NumericToFloat64(e.TotalAmount),
		ExpenseDate:      e.ExpenseDate.Time,
		DueDate:          db.DateToTimePtr(e.DueDate),
		CategoryID:       db.UUIDToStringPtr(e.CategoryID),
		SplitType:        e.SplitType,
		IsInstallment:    e.IsInstallment,
		InstallmentCount: db.Int2ToInt16Ptr(e.InstallmentCount),
		CreatedAt:        e.CreatedAt.Time,
		UpdatedAt:        e.UpdatedAt.Time,
	}
}

func toExpenses(expenses []db.Expense) []Expense {
	result := make([]Expense, len(expenses))
	for i, e := range expenses {
		result[i] = toExpense(e)
	}
	return result
}
