package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type ExpenseSplitRepo struct {
	q *db.Queries
}

func NewExpenseSplitRepo(q *db.Queries) *ExpenseSplitRepo {
	return &ExpenseSplitRepo{q: q}
}

func (r *ExpenseSplitRepo) GetByID(ctx context.Context, id string) (domain.ExpenseSplit, error) {
	es, err := r.q.GetExpenseSplitByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.ExpenseSplit{}, err
	}
	return domainExpenseSplit(es), nil
}

func (r *ExpenseSplitRepo) ListByExpense(ctx context.Context, expenseID string) ([]domain.ExpenseSplitWithUser, error) {
	rows, err := r.q.ListExpenseSplitsByExpense(ctx, uuidFromString(expenseID))
	if err != nil {
		return nil, err
	}
	return domainExpenseSplitsWithUser(rows), nil
}

func (r *ExpenseSplitRepo) ListByUser(ctx context.Context, userID, groupID string) ([]domain.ExpenseSplit, error) {
	rows, err := r.q.ListExpenseSplitsByUser(ctx, db.ListExpenseSplitsByUserParams{
		UserID:  uuidFromString(userID),
		GroupID: uuidFromString(groupID),
	})
	if err != nil {
		return nil, err
	}
	result := make([]domain.ExpenseSplit, len(rows))
	for i, row := range rows {
		result[i] = domain.ExpenseSplit{
			ID:        uuidToString(row.ID),
			ExpenseID: uuidToString(row.ExpenseID),
			UserID:    uuidToString(row.UserID),
			Amount:    numericToFloat64(row.Amount),
			IsPaid:    row.IsPaid,
			PaidAt:    timestamptzToTimePtr(row.PaidAt),
			CreatedAt: row.CreatedAt.Time,
		}
	}
	return result, nil
}

func (r *ExpenseSplitRepo) Create(ctx context.Context, expenseID, userID string, amount float64) (domain.ExpenseSplit, error) {
	es, err := r.q.CreateExpenseSplit(ctx, db.CreateExpenseSplitParams{
		ExpenseID: uuidFromString(expenseID),
		UserID:    uuidFromString(userID),
		Amount:    numericFromFloat64(amount),
	})
	if err != nil {
		return domain.ExpenseSplit{}, err
	}
	return domainExpenseSplit(es), nil
}

func (r *ExpenseSplitRepo) CreateMany(ctx context.Context, expenseID string, splits []domain.CreateExpenseSplitParams) (int64, error) {
	params := make([]db.CreateExpenseSplitsParams, len(splits))
	for i, s := range splits {
		params[i] = db.CreateExpenseSplitsParams{
			ExpenseID: uuidFromString(expenseID),
			UserID:    uuidFromString(s.UserID),
			Amount:    numericFromFloat64(s.Amount),
		}
	}
	return r.q.CreateExpenseSplits(ctx, params)
}

func (r *ExpenseSplitRepo) MarkAsPaid(ctx context.Context, id string) error {
	return r.q.MarkExpenseSplitAsPaid(ctx, uuidFromString(id))
}

func (r *ExpenseSplitRepo) DeleteByExpense(ctx context.Context, expenseID string) error {
	return r.q.DeleteExpenseSplitsByExpense(ctx, uuidFromString(expenseID))
}

func (r *ExpenseSplitRepo) GetUserBalance(ctx context.Context, userID, groupID string) (domain.UserBalance, error) {
	row, err := r.q.GetUserBalanceInGroup(ctx, db.GetUserBalanceInGroupParams{
		UserID:  uuidFromString(userID),
		GroupID: uuidFromString(groupID),
	})
	if err != nil {
		return domain.UserBalance{}, err
	}
	return domainUserBalance(row), nil
}
