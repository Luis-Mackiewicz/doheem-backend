package expense

import (
	"context"
	"time"

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

func (r *ExpenseRepo) ListByParent(ctx context.Context, parentID string) ([]Expense, error) {
	expenses, err := r.q.ListExpensesByParent(ctx, db.UUIDFromString(parentID))
	if err != nil {
		return nil, err
	}
	return toExpenses(expenses), nil
}

func (r *ExpenseRepo) Create(ctx context.Context, params CreateExpenseParams) (Expense, error) {
	e, err := r.q.CreateExpense(ctx, db.CreateExpenseParams{
		GroupID:          db.UUIDFromString(params.GroupID),
		Description:      params.Description,
		Amount:           db.NumericFromFloat64(params.Amount),
		CategoryID:       db.UUIDFromString(params.CategoryID),
		CompetenceDate:   db.DateFromTime(params.CompetenceDate),
		DueDate:          db.DateFromTime(params.DueDate),
		PaidBy:           db.UUIDFromString(params.PaidBy),
		SplitMode:        params.SplitMode,
		Installments:     params.Installments,
		FirstDueDate:     dateFromTimePtr(params.FirstDueDate),
		IsFixed:          params.IsFixed,
		ParentExpenseID:  db.UUIDFromStringPtr(params.ParentExpenseID),
		InstallmentIndex: int4FromInt32Ptr(params.InstallmentIndex),
		InstallmentTotal: int4FromInt32Ptr(params.InstallmentTotal),
	})
	if err != nil {
		return Expense{}, err
	}
	return toExpense(e), nil
}

func (r *ExpenseRepo) Update(ctx context.Context, id string, params UpdateExpenseParams) (Expense, error) {
	e, err := r.q.UpdateExpense(ctx, db.UpdateExpenseParams{
		ID:             db.UUIDFromString(id),
		Description:    deptrStr(params.Description),
		Amount:         deptrNumeric(params.Amount),
		CompetenceDate: deptrDate(params.CompetenceDate),
		DueDate:        deptrDate(params.DueDate),
		CategoryID:     db.UUIDFromStringPtrOrZero(params.CategoryID),
		SplitMode:      deptrStr(params.SplitMode),
	})
	if err != nil {
		return Expense{}, err
	}
	return toExpense(e), nil
}

func (r *ExpenseRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeleteExpense(ctx, db.UUIDFromString(id))
}

func (r *ExpenseRepo) DeleteByParent(ctx context.Context, parentID string) error {
	return r.q.DeleteExpensesByParent(ctx, db.UUIDFromString(parentID))
}

func (r *ExpenseRepo) GetTotalByGroup(ctx context.Context, groupID string) (float64, error) {
	val, err := r.q.GetTotalExpensesByGroup(ctx, db.UUIDFromString(groupID))
	if err != nil {
		return 0, err
	}
	return db.NumericToFloat64(val), nil
}

func toExpense(e db.Expense) Expense {
	var installmentIndex *int32
	if e.InstallmentIndex.Valid {
		installmentIndex = &e.InstallmentIndex.Int32
	}
	var installmentTotal *int32
	if e.InstallmentTotal.Valid {
		installmentTotal = &e.InstallmentTotal.Int32
	}
	return Expense{
		ID:               db.UUIDToString(e.ID),
		GroupID:          db.UUIDToString(e.GroupID),
		Description:      e.Description,
		Amount:           db.NumericToFloat64(e.Amount),
		CategoryID:       db.UUIDToString(e.CategoryID),
		CompetenceDate:   e.CompetenceDate.Time,
		DueDate:          e.DueDate.Time,
		PaidBy:           db.UUIDToString(e.PaidBy),
		SplitMode:        e.SplitMode,
		Installments:     e.Installments,
		FirstDueDate:     db.DateToTimePtr(e.FirstDueDate),
		IsFixed:          e.IsFixed,
		ParentExpenseID:  db.UUIDToStringPtr(e.ParentExpenseID),
		InstallmentIndex: installmentIndex,
		InstallmentTotal: installmentTotal,
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

func deptrDate(t *time.Time) pgtype.Date {
	if t != nil {
		return db.DateFromTime(*t)
	}
	return pgtype.Date{}
}

func dateFromTimePtr(t *time.Time) pgtype.Date {
	if t != nil {
		return db.DateFromTime(*t)
	}
	return pgtype.Date{}
}

func int4FromInt32Ptr(i *int32) pgtype.Int4 {
	if i != nil {
		return pgtype.Int4{Int32: *i, Valid: true}
	}
	return pgtype.Int4{}
}
