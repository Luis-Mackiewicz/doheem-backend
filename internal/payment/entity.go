package payment

import (
	"context"
	"errors"
	"time"
)

type Payment struct {
	ID          string
	GroupID     string
	PayerID     string
	ReceiverID  string
	Amount      float64
	PaymentDate time.Time
	Status      string
	Notes       *string
	CreatedAt   time.Time
	ConfirmedAt *time.Time
	CancelledAt *time.Time
}

type PaymentWithUsers struct {
	Payment
	PayerName   string
	ReceiverName string
}

type PaymentWithGroup struct {
	Payment
	GroupName string
}

type CreatePaymentParams struct {
	GroupID     string
	PayerID     string
	ReceiverID  string
	Amount      float64
	PaymentDate time.Time
	Notes       *string
}

type PaymentAttachment struct {
	ID         string
	PaymentID  string
	FilePath   string
	FileType   string
	FileSize   int32
	UploadedAt time.Time
}

type PaymentRepository interface {
	GetByID(ctx context.Context, id string) (Payment, error)
	ListByGroup(ctx context.Context, groupID string) ([]PaymentWithUsers, error)
	ListByUser(ctx context.Context, userID string) ([]PaymentWithGroup, error)
	ListPendingByUser(ctx context.Context, userID string) ([]Payment, error)
	Create(ctx context.Context, params CreatePaymentParams) (Payment, error)
	Confirm(ctx context.Context, id string) error
	Cancel(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

var (
	ErrPaymentNotFound           = errors.New("payment not found")
	ErrPaymentAlreadyConfirmed   = errors.New("payment already confirmed")
	ErrPaymentAlreadyCancelled   = errors.New("payment already cancelled")
)

type PaymentAttachmentRepository interface {
	GetByID(ctx context.Context, id string) (PaymentAttachment, error)
	ListByPayment(ctx context.Context, paymentID string) ([]PaymentAttachment, error)
	Create(ctx context.Context, paymentID, filePath, fileType string, fileSize int32) (PaymentAttachment, error)
	Delete(ctx context.Context, id string) error
}
