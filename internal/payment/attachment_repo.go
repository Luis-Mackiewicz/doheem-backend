package payment

import (
	"context"

	"doheem-backend/internal/db"
)

type PaymentAttachmentRepo struct {
	q *db.Queries
}

func NewPaymentAttachmentRepo(q *db.Queries) *PaymentAttachmentRepo {
	return &PaymentAttachmentRepo{q: q}
}

func (r *PaymentAttachmentRepo) GetByID(ctx context.Context, id string) (PaymentAttachment, error) {
	pa, err := r.q.GetPaymentAttachmentByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return PaymentAttachment{}, err
	}
	return toAttachment(pa), nil
}

func (r *PaymentAttachmentRepo) ListByPayment(ctx context.Context, paymentID string) ([]PaymentAttachment, error) {
	attachments, err := r.q.ListPaymentAttachmentsByPayment(ctx, db.UUIDFromString(paymentID))
	if err != nil {
		return nil, err
	}
	return toAttachments(attachments), nil
}

func (r *PaymentAttachmentRepo) Create(ctx context.Context, paymentID, filePath, fileType string, fileSize int32) (PaymentAttachment, error) {
	pa, err := r.q.CreatePaymentAttachment(ctx, db.CreatePaymentAttachmentParams{
		PaymentID: db.UUIDFromString(paymentID),
		FilePath:  filePath,
		FileType:  fileType,
		FileSize:  fileSize,
	})
	if err != nil {
		return PaymentAttachment{}, err
	}
	return toAttachment(pa), nil
}

func (r *PaymentAttachmentRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeletePaymentAttachment(ctx, db.UUIDFromString(id))
}

func toAttachment(pa db.PaymentAttachment) PaymentAttachment {
	return PaymentAttachment{
		ID:         db.UUIDToString(pa.ID),
		PaymentID:  db.UUIDToString(pa.PaymentID),
		FilePath:   pa.FilePath,
		FileType:   pa.FileType,
		FileSize:   pa.FileSize,
		UploadedAt: pa.UploadedAt.Time,
	}
}

func toAttachments(attachments []db.PaymentAttachment) []PaymentAttachment {
	result := make([]PaymentAttachment, len(attachments))
	for i, a := range attachments {
		result[i] = toAttachment(a)
	}
	return result
}
