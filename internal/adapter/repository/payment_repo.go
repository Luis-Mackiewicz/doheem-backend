package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
)

type PaymentRepo struct {
	q *db.Queries
}

func NewPaymentRepo(q *db.Queries) *PaymentRepo {
	return &PaymentRepo{q: q}
}

func (r *PaymentRepo) GetByID(ctx context.Context, id string) (domain.Payment, error) {
	p, err := r.q.GetPaymentByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.Payment{}, err
	}
	return domainPayment(p), nil
}

func (r *PaymentRepo) ListByGroup(ctx context.Context, groupID string) ([]domain.PaymentWithUsers, error) {
	rows, err := r.q.ListPaymentsByGroup(ctx, uuidFromString(groupID))
	if err != nil {
		return nil, err
	}
	return domainPaymentsWithUsers(rows), nil
}

func (r *PaymentRepo) ListByUser(ctx context.Context, userID string) ([]domain.PaymentWithGroup, error) {
	rows, err := r.q.ListPaymentsByUser(ctx, uuidFromString(userID))
	if err != nil {
		return nil, err
	}
	return domainPaymentsWithGroup(rows), nil
}

func (r *PaymentRepo) ListPendingByUser(ctx context.Context, userID string) ([]domain.Payment, error) {
	payments, err := r.q.ListPendingPaymentsByUser(ctx, uuidFromString(userID))
	if err != nil {
		return nil, err
	}
	return domainPayments(payments), nil
}

func (r *PaymentRepo) Create(ctx context.Context, params domain.CreatePaymentParams) (domain.Payment, error) {
	var notes pgtype.Text
	if params.Notes != nil {
		notes = pgtype.Text{String: *params.Notes, Valid: true}
	}
	p, err := r.q.CreatePayment(ctx, db.CreatePaymentParams{
		GroupID:     uuidFromString(params.GroupID),
		PayerID:     uuidFromString(params.PayerID),
		ReceiverID:  uuidFromString(params.ReceiverID),
		Amount:      numericFromFloat64(params.Amount),
		PaymentDate: dateFromTime(params.PaymentDate),
		Notes:       notes,
	})
	if err != nil {
		return domain.Payment{}, err
	}
	return domainPayment(p), nil
}

func (r *PaymentRepo) Confirm(ctx context.Context, id string) error {
	return r.q.ConfirmPayment(ctx, uuidFromString(id))
}

func (r *PaymentRepo) Cancel(ctx context.Context, id string) error {
	return r.q.CancelPayment(ctx, uuidFromString(id))
}

func (r *PaymentRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeletePayment(ctx, uuidFromString(id))
}
