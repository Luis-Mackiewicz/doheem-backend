package expense

import (
	"context"

	"doheem-backend/internal/db"
)

type InstallmentRepo struct {
	q *db.Queries
}

func NewInstallmentRepo(q *db.Queries) *InstallmentRepo {
	return &InstallmentRepo{q: q}
}

func (r *InstallmentRepo) GetByID(ctx context.Context, id string) (Installment, error) {
	i, err := r.q.GetInstallmentByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return Installment{}, err
	}
	return toInstallment(i), nil
}

func (r *InstallmentRepo) ListByExpense(ctx context.Context, expenseID string) ([]Installment, error) {
	installments, err := r.q.ListInstallmentsByExpense(ctx, db.UUIDFromString(expenseID))
	if err != nil {
		return nil, err
	}
	return toInstallments(installments), nil
}

func (r *InstallmentRepo) Create(ctx context.Context, params CreateInstallmentParams) (Installment, error) {
	i, err := r.q.CreateInstallment(ctx, db.CreateInstallmentParams{
		ExpenseID:         db.UUIDFromString(""),
		InstallmentNumber: params.InstallmentNumber,
		Amount:            db.NumericFromFloat64(params.Amount),
		DueDate:           db.DateFromTime(params.DueDate),
	})
	if err != nil {
		return Installment{}, err
	}
	return toInstallment(i), nil
}

func (r *InstallmentRepo) CreateMany(ctx context.Context, expenseID string, installments []CreateInstallmentParams) (int64, error) {
	params := make([]db.CreateInstallmentsParams, len(installments))
	for i, inst := range installments {
		params[i] = db.CreateInstallmentsParams{
			ExpenseID:         db.UUIDFromString(expenseID),
			InstallmentNumber: inst.InstallmentNumber,
			Amount:            db.NumericFromFloat64(inst.Amount),
			DueDate:           db.DateFromTime(inst.DueDate),
		}
	}
	return r.q.CreateInstallments(ctx, params)
}

func (r *InstallmentRepo) MarkAsPaid(ctx context.Context, id string) error {
	return r.q.MarkInstallmentAsPaid(ctx, db.UUIDFromString(id))
}

func (r *InstallmentRepo) DeleteByExpense(ctx context.Context, expenseID string) error {
	return r.q.DeleteInstallmentsByExpense(ctx, db.UUIDFromString(expenseID))
}

func toInstallment(i db.Installment) Installment {
	return Installment{
		ID:                db.UUIDToString(i.ID),
		ExpenseID:         db.UUIDToString(i.ExpenseID),
		InstallmentNumber: i.InstallmentNumber,
		Amount:            db.NumericToFloat64(i.Amount),
		DueDate:           i.DueDate.Time,
		IsPaid:            i.IsPaid,
		PaidAt:            db.TimestamptzToTimePtr(i.PaidAt),
		CreatedAt:         i.CreatedAt.Time,
	}
}

func toInstallments(installments []db.Installment) []Installment {
	result := make([]Installment, len(installments))
	for i, inst := range installments {
		result[i] = toInstallment(inst)
	}
	return result
}
