package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type PaymentAttachmentRepo struct {
	q *db.Queries
}

func NewPaymentAttachmentRepo(q *db.Queries) *PaymentAttachmentRepo {
	return &PaymentAttachmentRepo{q: q}
}

func (r *PaymentAttachmentRepo) GetByID(ctx context.Context, id string) (domain.PaymentAttachment, error) {
	pa, err := r.q.GetPaymentAttachmentByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.PaymentAttachment{}, err
	}
	return domainPaymentAttachment(pa), nil
}

func (r *PaymentAttachmentRepo) ListByPayment(ctx context.Context, paymentID string) ([]domain.PaymentAttachment, error) {
	attachments, err := r.q.ListPaymentAttachmentsByPayment(ctx, uuidFromString(paymentID))
	if err != nil {
		return nil, err
	}
	return domainPaymentAttachments(attachments), nil
}

func (r *PaymentAttachmentRepo) Create(ctx context.Context, paymentID, filePath, fileType string, fileSize int32) (domain.PaymentAttachment, error) {
	pa, err := r.q.CreatePaymentAttachment(ctx, db.CreatePaymentAttachmentParams{
		PaymentID: uuidFromString(paymentID),
		FilePath:  filePath,
		FileType:  fileType,
		FileSize:  fileSize,
	})
	if err != nil {
		return domain.PaymentAttachment{}, err
	}
	return domainPaymentAttachment(pa), nil
}

func (r *PaymentAttachmentRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeletePaymentAttachment(ctx, uuidFromString(id))
}
