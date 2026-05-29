package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type InstallmentRepo struct {
	q *db.Queries
}

func NewInstallmentRepo(q *db.Queries) *InstallmentRepo {
	return &InstallmentRepo{q: q}
}

func (r *InstallmentRepo) GetByID(ctx context.Context, id string) (domain.Installment, error) {
	i, err := r.q.GetInstallmentByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.Installment{}, err
	}
	return domainInstallment(i), nil
}

func (r *InstallmentRepo) ListByExpense(ctx context.Context, expenseID string) ([]domain.Installment, error) {
	installments, err := r.q.ListInstallmentsByExpense(ctx, uuidFromString(expenseID))
	if err != nil {
		return nil, err
	}
	return domainInstallments(installments), nil
}

func (r *InstallmentRepo) Create(ctx context.Context, params domain.CreateInstallmentParams) (domain.Installment, error) {
	i, err := r.q.CreateInstallment(ctx, db.CreateInstallmentParams{
		ExpenseID:         uuidFromString(""),
		InstallmentNumber: params.InstallmentNumber,
		Amount:            numericFromFloat64(params.Amount),
		DueDate:           dateFromTime(params.DueDate),
	})
	if err != nil {
		return domain.Installment{}, err
	}
	return domainInstallment(i), nil
}

func (r *InstallmentRepo) CreateMany(ctx context.Context, expenseID string, installments []domain.CreateInstallmentParams) (int64, error) {
	params := make([]db.CreateInstallmentsParams, len(installments))
	for i, inst := range installments {
		params[i] = db.CreateInstallmentsParams{
			ExpenseID:         uuidFromString(expenseID),
			InstallmentNumber: inst.InstallmentNumber,
			Amount:            numericFromFloat64(inst.Amount),
			DueDate:           dateFromTime(inst.DueDate),
		}
	}
	return r.q.CreateInstallments(ctx, params)
}

func (r *InstallmentRepo) MarkAsPaid(ctx context.Context, id string) error {
	return r.q.MarkInstallmentAsPaid(ctx, uuidFromString(id))
}

func (r *InstallmentRepo) DeleteByExpense(ctx context.Context, expenseID string) error {
	return r.q.DeleteInstallmentsByExpense(ctx, uuidFromString(expenseID))
}
