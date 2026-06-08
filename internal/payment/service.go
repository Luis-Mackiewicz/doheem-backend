package payment

import (
	"context"
)

type PaymentService struct {
	paymentRepo         PaymentRepository
	paymentAttachmentRepo PaymentAttachmentRepository
}

func NewPaymentService(paymentRepo PaymentRepository, paymentAttachmentRepo PaymentAttachmentRepository) *PaymentService {
	return &PaymentService{
		paymentRepo:         paymentRepo,
		paymentAttachmentRepo: paymentAttachmentRepo,
	}
}

func (s *PaymentService) Create(ctx context.Context, params CreatePaymentParams) (Payment, error) {
	return s.paymentRepo.Create(ctx, params)
}

func (s *PaymentService) GetByID(ctx context.Context, id string) (Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return Payment{}, ErrPaymentNotFound
	}
	return payment, nil
}

func (s *PaymentService) ListByGroup(ctx context.Context, groupID string) ([]PaymentWithUsers, error) {
	return s.paymentRepo.ListByGroup(ctx, groupID)
}

func (s *PaymentService) ListByUser(ctx context.Context, userID string) ([]PaymentWithGroup, error) {
	return s.paymentRepo.ListByUser(ctx, userID)
}

func (s *PaymentService) ListPendingByUser(ctx context.Context, userID string) ([]Payment, error) {
	return s.paymentRepo.ListPendingByUser(ctx, userID)
}

func (s *PaymentService) Confirm(ctx context.Context, id string) error {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return ErrPaymentNotFound
	}
	if payment.Status == "confirmed" {
		return ErrPaymentAlreadyConfirmed
	}
	if payment.Status == "cancelled" {
		return ErrPaymentAlreadyCancelled
	}
	if err := s.paymentRepo.Confirm(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *PaymentService) Cancel(ctx context.Context, id string) error {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return ErrPaymentNotFound
	}
	if payment.Status == "cancelled" {
		return ErrPaymentAlreadyCancelled
	}
	return s.paymentRepo.Cancel(ctx, id)
}

func (s *PaymentService) Delete(ctx context.Context, id string) error {
	return s.paymentRepo.Delete(ctx, id)
}

func (s *PaymentService) AddAttachment(ctx context.Context, paymentID, filePath, fileType string, fileSize int32) (PaymentAttachment, error) {
	return s.paymentAttachmentRepo.Create(ctx, paymentID, filePath, fileType, fileSize)
}

func (s *PaymentService) ListAttachments(ctx context.Context, paymentID string) ([]PaymentAttachment, error) {
	return s.paymentAttachmentRepo.ListByPayment(ctx, paymentID)
}

func (s *PaymentService) DeleteAttachment(ctx context.Context, attachmentID string) error {
	return s.paymentAttachmentRepo.Delete(ctx, attachmentID)
}
