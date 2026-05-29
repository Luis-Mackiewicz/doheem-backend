package repository

import (
	"context"
	"testing"
	"time"

	"doheem-backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentRepo_CreateAndGetByID(t *testing.T) {
	q := newTxQueries(t)
	paymentRepo := NewPaymentRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group := createTestGroup(t, q)

	payment, err := paymentRepo.Create(ctx, domain.CreatePaymentParams{
		GroupID:     group.ID,
		PayerID:     user.ID,
		ReceiverID:  user.ID,
		Amount:      150.00,
		PaymentDate: time.Now(),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, payment.ID)
	assert.Equal(t, "pending", payment.Status)
	assert.Equal(t, 150.00, payment.Amount)

	got, err := paymentRepo.GetByID(ctx, payment.ID)
	require.NoError(t, err)
	assert.Equal(t, payment.ID, got.ID)
}

func TestPaymentRepo_Confirm(t *testing.T) {
	q := newTxQueries(t)
	paymentRepo := NewPaymentRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group := createTestGroup(t, q)

	payment, err := paymentRepo.Create(ctx, domain.CreatePaymentParams{
		GroupID:     group.ID,
		PayerID:     user.ID,
		ReceiverID:  user.ID,
		Amount:      100.00,
		PaymentDate: time.Now(),
	})
	require.NoError(t, err)

	err = paymentRepo.Confirm(ctx, payment.ID)
	require.NoError(t, err)

	got, err := paymentRepo.GetByID(ctx, payment.ID)
	require.NoError(t, err)
	assert.Equal(t, "confirmed", got.Status)
	assert.NotNil(t, got.ConfirmedAt)
}

func TestPaymentRepo_Cancel(t *testing.T) {
	q := newTxQueries(t)
	paymentRepo := NewPaymentRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group := createTestGroup(t, q)

	payment, err := paymentRepo.Create(ctx, domain.CreatePaymentParams{
		GroupID:     group.ID,
		PayerID:     user.ID,
		ReceiverID:  user.ID,
		Amount:      100.00,
		PaymentDate: time.Now(),
	})
	require.NoError(t, err)

	err = paymentRepo.Cancel(ctx, payment.ID)
	require.NoError(t, err)

	got, err := paymentRepo.GetByID(ctx, payment.ID)
	require.NoError(t, err)
	assert.Equal(t, "cancelled", got.Status)
	assert.NotNil(t, got.CancelledAt)
}

func TestPaymentRepo_ListByGroup(t *testing.T) {
	q := newTxQueries(t)
	paymentRepo := NewPaymentRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group := createTestGroup(t, q)

	paymentRepo.Create(ctx, domain.CreatePaymentParams{
		GroupID: group.ID, PayerID: user.ID, ReceiverID: user.ID,
		Amount: 50, PaymentDate: time.Now(),
	})
	paymentRepo.Create(ctx, domain.CreatePaymentParams{
		GroupID: group.ID, PayerID: user.ID, ReceiverID: user.ID,
		Amount: 75, PaymentDate: time.Now(),
	})

	payments, err := paymentRepo.ListByGroup(ctx, group.ID)
	require.NoError(t, err)
	assert.Len(t, payments, 2)
}

func TestPaymentRepo_Delete(t *testing.T) {
	q := newTxQueries(t)
	paymentRepo := NewPaymentRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group := createTestGroup(t, q)

	payment, err := paymentRepo.Create(ctx, domain.CreatePaymentParams{
		GroupID: group.ID, PayerID: user.ID, ReceiverID: user.ID,
		Amount: 200, PaymentDate: time.Now(),
	})
	require.NoError(t, err)

	err = paymentRepo.Delete(ctx, payment.ID)
	require.NoError(t, err)

	_, err = paymentRepo.GetByID(ctx, payment.ID)
	assert.Error(t, err)
}
