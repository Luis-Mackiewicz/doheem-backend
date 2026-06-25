package expense

import (
	"context"

	"doheem-backend/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type ExpenseSplitRepo struct {
	q *db.Queries
}

func NewExpenseSplitRepo(q *db.Queries) *ExpenseSplitRepo {
	return &ExpenseSplitRepo{q: q}
}

func (r *ExpenseSplitRepo) GetByID(ctx context.Context, id string) (ExpenseSplit, error) {
	es, err := r.q.GetExpenseSplitByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return ExpenseSplit{}, err
	}
	return toExpenseSplit(es), nil
}

func (r *ExpenseSplitRepo) ListByExpense(ctx context.Context, expenseID string) ([]ExpenseSplitWithUser, error) {
	rows, err := r.q.ListExpenseSplitsByExpense(ctx, db.UUIDFromString(expenseID))
	if err != nil {
		return nil, err
	}
	return toExpenseSplitsWithUser(rows), nil
}

func (r *ExpenseSplitRepo) ListByExpenseIDs(ctx context.Context, expenseIDs []string) ([]ExpenseSplitWithUser, error) {
	uuids := make([]pgtype.UUID, len(expenseIDs))
	for i, id := range expenseIDs {
		uuids[i] = db.UUIDFromString(id)
	}
	if len(uuids) == 0 {
		return nil, nil
	}
	rows, err := r.q.ListExpenseSplitsByExpenseIDs(ctx, uuids)
	if err != nil {
		return nil, err
	}
	return toExpenseSplitsWithUserFromIDsRow(rows), nil
}

func (r *ExpenseSplitRepo) ListByUser(ctx context.Context, userID, groupID string) ([]ExpenseSplit, error) {
	rows, err := r.q.ListExpenseSplitsByUser(ctx, db.ListExpenseSplitsByUserParams{
		UserID:  db.UUIDFromString(userID),
		GroupID: db.UUIDFromString(groupID),
	})
	if err != nil {
		return nil, err
	}
	result := make([]ExpenseSplit, len(rows))
	for i, row := range rows {
		result[i] = ExpenseSplit{
			ID:              db.UUIDToString(row.ID),
			ExpenseID:       db.UUIDToString(row.ExpenseID),
			UserID:          db.UUIDToString(row.UserID),
			Amount:          db.NumericToDecimal(row.Amount),
			IsPaid:          row.IsPaid,
			PaidAt:          db.TimestamptzToTimePtr(row.PaidAt),
			CreatedAt:       row.CreatedAt.Time,
			ReceiptData:     db.TextToStringPtr(row.ReceiptData),
			ReceiptType:     db.TextToStringPtr(row.ReceiptType),
			ReceiptFileName: db.TextToStringPtr(row.ReceiptFileName),
		}
	}
	return result, nil
}

func (r *ExpenseSplitRepo) Create(ctx context.Context, expenseID, userID string, amount decimal.Decimal) (ExpenseSplit, error) {
	es, err := r.q.CreateExpenseSplit(ctx, db.CreateExpenseSplitParams{
		ExpenseID: db.UUIDFromString(expenseID),
		UserID:    db.UUIDFromString(userID),
		Amount:    db.NumericFromDecimal(amount),
	})
	if err != nil {
		return ExpenseSplit{}, err
	}
	return toExpenseSplit(es), nil
}

func (r *ExpenseSplitRepo) CreateMany(ctx context.Context, expenseID string, splits []CreateExpenseSplitParams) (int64, error) {
	params := make([]db.CreateExpenseSplitsParams, len(splits))
	for i, s := range splits {
		params[i] = db.CreateExpenseSplitsParams{
			ExpenseID: db.UUIDFromString(expenseID),
			UserID:    db.UUIDFromString(s.UserID),
			Amount:    db.NumericFromDecimal(s.Amount),
		}
	}
	return r.q.CreateExpenseSplits(ctx, params)
}

func (r *ExpenseSplitRepo) MarkAsPaid(ctx context.Context, id string, receiptData, receiptType, receiptFileName *string) error {
	return r.q.MarkExpenseSplitAsPaid(ctx, db.MarkExpenseSplitAsPaidParams{
		ID:              db.UUIDFromString(id),
		ReceiptData:     db.TextFromStringPtr(receiptData),
		ReceiptType:     db.TextFromStringPtr(receiptType),
		ReceiptFileName: db.TextFromStringPtr(receiptFileName),
	})
}

func (r *ExpenseSplitRepo) MarkAsPaidByExpenseAndUserIDs(ctx context.Context, expenseID string, userIDs []string) error {
	uuids := make([]pgtype.UUID, len(userIDs))
	for i, uid := range userIDs {
		uuids[i] = db.UUIDFromString(uid)
	}
	return r.q.MarkExpenseSplitsAsPaidByExpense(ctx, db.MarkExpenseSplitsAsPaidByExpenseParams{
		ExpenseID: db.UUIDFromString(expenseID),
		UserIDs:   uuids,
	})
}

func (r *ExpenseSplitRepo) HasPaidSplits(ctx context.Context, expenseID string) (bool, error) {
	return r.q.HasExpensePaidSplits(ctx, db.UUIDFromString(expenseID))
}

func (r *ExpenseSplitRepo) DeleteByExpense(ctx context.Context, expenseID string) error {
	return r.q.DeleteExpenseSplitsByExpense(ctx, db.UUIDFromString(expenseID))
}

func (r *ExpenseSplitRepo) GetUserBalance(ctx context.Context, userID, groupID string) (UserBalance, error) {
	row, err := r.q.GetUserBalanceInGroup(ctx, db.GetUserBalanceInGroupParams{
		UserID:  db.UUIDFromString(userID),
		GroupID: db.UUIDFromString(groupID),
	})
	if err != nil {
		return UserBalance{}, err
	}
	return toUserBalance(row), nil
}

func toExpenseSplit(es db.ExpenseSplit) ExpenseSplit {
	return ExpenseSplit{
		ID:              db.UUIDToString(es.ID),
		ExpenseID:       db.UUIDToString(es.ExpenseID),
		UserID:          db.UUIDToString(es.UserID),
		Amount:          db.NumericToDecimal(es.Amount),
		IsPaid:          es.IsPaid,
		PaidAt:          db.TimestamptzToTimePtr(es.PaidAt),
		CreatedAt:       es.CreatedAt.Time,
		ReceiptData:     db.TextToStringPtr(es.ReceiptData),
		ReceiptType:     db.TextToStringPtr(es.ReceiptType),
		ReceiptFileName: db.TextToStringPtr(es.ReceiptFileName),
	}
}

func toExpenseSplitWithUser(row db.ListExpenseSplitsByExpenseRow) ExpenseSplitWithUser {
	return ExpenseSplitWithUser{
		ExpenseSplit: ExpenseSplit{
			ID:              db.UUIDToString(row.ID),
			ExpenseID:       db.UUIDToString(row.ExpenseID),
			UserID:          db.UUIDToString(row.UserID),
			Amount:          db.NumericToDecimal(row.Amount),
			IsPaid:          row.IsPaid,
			PaidAt:          db.TimestamptzToTimePtr(row.PaidAt),
			CreatedAt:       row.CreatedAt.Time,
			ReceiptData:     db.TextToStringPtr(row.ReceiptData),
			ReceiptType:     db.TextToStringPtr(row.ReceiptType),
			ReceiptFileName: db.TextToStringPtr(row.ReceiptFileName),
		},
		UserName:  row.UserName,
		UserEmail: row.UserEmail,
	}
}

func toExpenseSplitsWithUser(rows []db.ListExpenseSplitsByExpenseRow) []ExpenseSplitWithUser {
	result := make([]ExpenseSplitWithUser, len(rows))
	for i, r := range rows {
		result[i] = toExpenseSplitWithUser(r)
	}
	return result
}

func toExpenseSplitsWithUserFromIDsRow(rows []db.ListExpenseSplitsByExpenseIDsRow) []ExpenseSplitWithUser {
	result := make([]ExpenseSplitWithUser, len(rows))
	for i, r := range rows {
		result[i] = ExpenseSplitWithUser{
			ExpenseSplit: ExpenseSplit{
				ID:              db.UUIDToString(r.ID),
				ExpenseID:       db.UUIDToString(r.ExpenseID),
				UserID:          db.UUIDToString(r.UserID),
				Amount:          db.NumericToDecimal(r.Amount),
				IsPaid:          r.IsPaid,
				PaidAt:          db.TimestamptzToTimePtr(r.PaidAt),
				CreatedAt:       r.CreatedAt.Time,
				ReceiptData:     db.TextToStringPtr(r.ReceiptData),
				ReceiptType:     db.TextToStringPtr(r.ReceiptType),
				ReceiptFileName: db.TextToStringPtr(r.ReceiptFileName),
			},
			UserName:  r.UserName,
			UserEmail: r.UserEmail,
		}
	}
	return result
}

func toUserBalance(row db.GetUserBalanceInGroupRow) UserBalance {
	return UserBalance{
		TotalOwed: db.NumericToDecimal(row.TotalOwed),
		TotalPaid: db.NumericToDecimal(row.TotalPaid),
	}
}
