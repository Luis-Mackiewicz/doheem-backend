package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPaymentRepo struct {
	mock.Mock
}

func (m *mockPaymentRepo) GetByID(ctx context.Context, id string) (Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Payment), args.Error(1)
}

func (m *mockPaymentRepo) ListByGroup(ctx context.Context, groupID string) ([]PaymentWithUsers, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]PaymentWithUsers), args.Error(1)
}

func (m *mockPaymentRepo) ListByUser(ctx context.Context, userID string) ([]PaymentWithGroup, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]PaymentWithGroup), args.Error(1)
}

func (m *mockPaymentRepo) ListPendingByUser(ctx context.Context, userID string) ([]Payment, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]Payment), args.Error(1)
}

func (m *mockPaymentRepo) Create(ctx context.Context, params CreatePaymentParams) (Payment, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(Payment), args.Error(1)
}

func (m *mockPaymentRepo) Confirm(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockPaymentRepo) Cancel(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockPaymentRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockPaymentAttachmentRepo struct {
	mock.Mock
}

func (m *mockPaymentAttachmentRepo) GetByID(ctx context.Context, id string) (PaymentAttachment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(PaymentAttachment), args.Error(1)
}

func (m *mockPaymentAttachmentRepo) ListByPayment(ctx context.Context, paymentID string) ([]PaymentAttachment, error) {
	args := m.Called(ctx, paymentID)
	return args.Get(0).([]PaymentAttachment), args.Error(1)
}

func (m *mockPaymentAttachmentRepo) Create(ctx context.Context, paymentID, filePath, fileType string, fileSize int32) (PaymentAttachment, error) {
	args := m.Called(ctx, paymentID, filePath, fileType, fileSize)
	return args.Get(0).(PaymentAttachment), args.Error(1)
}

func (m *mockPaymentAttachmentRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestPaymentService_Confirm_Success(t *testing.T) {
	mockPayment := new(mockPaymentRepo)
	svc := NewPaymentService(mockPayment, new(mockPaymentAttachmentRepo))
	ctx := context.Background()

	mockPayment.On("GetByID", ctx, "p1").Return(Payment{ID: "p1", Status: "pending", GroupID: "g1", PayerID: "u1", ReceiverID: "u2", Amount: 100}, nil)
	mockPayment.On("Confirm", ctx, "p1").Return(nil)

	err := svc.Confirm(ctx, "p1")

	assert.NoError(t, err)
	mockPayment.AssertExpectations(t)
}

func TestPaymentService_Confirm_AlreadyConfirmed(t *testing.T) {
	mockPayment := new(mockPaymentRepo)
	svc := NewPaymentService(mockPayment, new(mockPaymentAttachmentRepo))
	ctx := context.Background()

	mockPayment.On("GetByID", ctx, "p1").Return(Payment{ID: "p1", Status: "confirmed"}, nil)

	err := svc.Confirm(ctx, "p1")

	assert.ErrorIs(t, err, ErrPaymentAlreadyConfirmed)
	mockPayment.AssertNotCalled(t, "Confirm")
}

func TestPaymentService_Confirm_AlreadyCancelled(t *testing.T) {
	mockPayment := new(mockPaymentRepo)
	svc := NewPaymentService(mockPayment, new(mockPaymentAttachmentRepo))
	ctx := context.Background()

	mockPayment.On("GetByID", ctx, "p1").Return(Payment{ID: "p1", Status: "cancelled"}, nil)

	err := svc.Confirm(ctx, "p1")

	assert.ErrorIs(t, err, ErrPaymentAlreadyCancelled)
	mockPayment.AssertNotCalled(t, "Confirm")
}

func TestPaymentService_Cancel_Success(t *testing.T) {
	mockPayment := new(mockPaymentRepo)
	svc := NewPaymentService(mockPayment, new(mockPaymentAttachmentRepo))
	ctx := context.Background()

	mockPayment.On("GetByID", ctx, "p1").Return(Payment{ID: "p1", Status: "pending"}, nil)
	mockPayment.On("Cancel", ctx, "p1").Return(nil)

	err := svc.Cancel(ctx, "p1")

	assert.NoError(t, err)
	mockPayment.AssertExpectations(t)
}

func TestPaymentService_Cancel_AlreadyCancelled(t *testing.T) {
	mockPayment := new(mockPaymentRepo)
	svc := NewPaymentService(mockPayment, new(mockPaymentAttachmentRepo))
	ctx := context.Background()

	mockPayment.On("GetByID", ctx, "p1").Return(Payment{ID: "p1", Status: "cancelled"}, nil)

	err := svc.Cancel(ctx, "p1")

	assert.ErrorIs(t, err, ErrPaymentAlreadyCancelled)
	mockPayment.AssertNotCalled(t, "Cancel")
}

func TestPaymentService_Confirm_NotFound(t *testing.T) {
	mockPayment := new(mockPaymentRepo)
	svc := NewPaymentService(mockPayment, new(mockPaymentAttachmentRepo))
	ctx := context.Background()

	mockPayment.On("GetByID", ctx, "999").Return(Payment{}, assert.AnError)

	err := svc.Confirm(ctx, "999")

	assert.ErrorIs(t, err, ErrPaymentNotFound)
}

func TestPaymentService_GetByID_NotFound(t *testing.T) {
	mockPayment := new(mockPaymentRepo)
	svc := NewPaymentService(mockPayment, new(mockPaymentAttachmentRepo))
	ctx := context.Background()

	mockPayment.On("GetByID", ctx, "999").Return(Payment{}, assert.AnError)

	_, err := svc.GetByID(ctx, "999")

	assert.ErrorIs(t, err, ErrPaymentNotFound)
}
