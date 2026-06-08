package payment

import (
	"context"

	"doheem-backend/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type PaymentRepo struct {
	q *db.Queries
}

func NewPaymentRepo(q *db.Queries) *PaymentRepo {
	return &PaymentRepo{q: q}
}

func (r *PaymentRepo) GetByID(ctx context.Context, id string) (Payment, error) {
	p, err := r.q.GetPaymentByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return Payment{}, err
	}
	return toPayment(p), nil
}

func (r *PaymentRepo) ListByGroup(ctx context.Context, groupID string) ([]PaymentWithUsers, error) {
	rows, err := r.q.ListPaymentsByGroup(ctx, db.UUIDFromString(groupID))
	if err != nil {
		return nil, err
	}
	return toPaymentsWithUsers(rows), nil
}

func (r *PaymentRepo) ListByUser(ctx context.Context, userID string) ([]PaymentWithGroup, error) {
	rows, err := r.q.ListPaymentsByUser(ctx, db.UUIDFromString(userID))
	if err != nil {
		return nil, err
	}
	return toPaymentsWithGroup(rows), nil
}

func (r *PaymentRepo) ListPendingByUser(ctx context.Context, userID string) ([]Payment, error) {
	payments, err := r.q.ListPendingPaymentsByUser(ctx, db.UUIDFromString(userID))
	if err != nil {
		return nil, err
	}
	return toPayments(payments), nil
}

func (r *PaymentRepo) Create(ctx context.Context, params CreatePaymentParams) (Payment, error) {
	var notes pgtype.Text
	if params.Notes != nil {
		notes = pgtype.Text{String: *params.Notes, Valid: true}
	}
	p, err := r.q.CreatePayment(ctx, db.CreatePaymentParams{
		GroupID:     db.UUIDFromString(params.GroupID),
		PayerID:     db.UUIDFromString(params.PayerID),
		ReceiverID:  db.UUIDFromString(params.ReceiverID),
		Amount:      db.NumericFromFloat64(params.Amount),
		PaymentDate: db.DateFromTime(params.PaymentDate),
		Notes:       notes,
	})
	if err != nil {
		return Payment{}, err
	}
	return toPayment(p), nil
}

func (r *PaymentRepo) Confirm(ctx context.Context, id string) error {
	return r.q.ConfirmPayment(ctx, db.UUIDFromString(id))
}

func (r *PaymentRepo) Cancel(ctx context.Context, id string) error {
	return r.q.CancelPayment(ctx, db.UUIDFromString(id))
}

func (r *PaymentRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeletePayment(ctx, db.UUIDFromString(id))
}

func toPayment(p db.Payment) Payment {
	return Payment{
		ID:          db.UUIDToString(p.ID),
		GroupID:     db.UUIDToString(p.GroupID),
		PayerID:     db.UUIDToString(p.PayerID),
		ReceiverID:  db.UUIDToString(p.ReceiverID),
		Amount:      db.NumericToFloat64(p.Amount),
		PaymentDate: p.PaymentDate.Time,
		Status:      p.Status,
		Notes:       db.TextToStringPtr(p.Notes),
		CreatedAt:   p.CreatedAt.Time,
		ConfirmedAt: db.TimestamptzToTimePtr(p.ConfirmedAt),
		CancelledAt: db.TimestamptzToTimePtr(p.CancelledAt),
	}
}

func toPayments(payments []db.Payment) []Payment {
	result := make([]Payment, len(payments))
	for i, p := range payments {
		result[i] = toPayment(p)
	}
	return result
}

func toPaymentWithUsers(row db.ListPaymentsByGroupRow) PaymentWithUsers {
	return PaymentWithUsers{
		Payment: Payment{
			ID:          db.UUIDToString(row.ID),
			GroupID:     db.UUIDToString(row.GroupID),
			PayerID:     db.UUIDToString(row.PayerID),
			ReceiverID:  db.UUIDToString(row.ReceiverID),
			Amount:      db.NumericToFloat64(row.Amount),
			PaymentDate: row.PaymentDate.Time,
			Status:      row.Status,
			Notes:       db.TextToStringPtr(row.Notes),
			CreatedAt:   row.CreatedAt.Time,
			ConfirmedAt: db.TimestamptzToTimePtr(row.ConfirmedAt),
			CancelledAt: db.TimestamptzToTimePtr(row.CancelledAt),
		},
		PayerName:   row.PayerName,
		ReceiverName: row.ReceiverName,
	}
}

func toPaymentsWithUsers(rows []db.ListPaymentsByGroupRow) []PaymentWithUsers {
	result := make([]PaymentWithUsers, len(rows))
	for i, r := range rows {
		result[i] = toPaymentWithUsers(r)
	}
	return result
}

func toPaymentWithGroup(row db.ListPaymentsByUserRow) PaymentWithGroup {
	return PaymentWithGroup{
		Payment: Payment{
			ID:          db.UUIDToString(row.ID),
			GroupID:     db.UUIDToString(row.GroupID),
			PayerID:     db.UUIDToString(row.PayerID),
			ReceiverID:  db.UUIDToString(row.ReceiverID),
			Amount:      db.NumericToFloat64(row.Amount),
			PaymentDate: row.PaymentDate.Time,
			Status:      row.Status,
			Notes:       db.TextToStringPtr(row.Notes),
			CreatedAt:   row.CreatedAt.Time,
			ConfirmedAt: db.TimestamptzToTimePtr(row.ConfirmedAt),
			CancelledAt: db.TimestamptzToTimePtr(row.CancelledAt),
		},
		GroupName: row.GroupName,
	}
}

func toPaymentsWithGroup(rows []db.ListPaymentsByUserRow) []PaymentWithGroup {
	result := make([]PaymentWithGroup, len(rows))
	for i, r := range rows {
		result[i] = toPaymentWithGroup(r)
	}
	return result
}
