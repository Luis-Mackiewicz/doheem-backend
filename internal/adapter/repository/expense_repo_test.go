package repository

import (
	"context"
	"testing"
	"time"

	"doheem-backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpenseRepo_CreateAndGetByID(t *testing.T) {
	q := newTxQueries(t)
	expenseRepo := NewExpenseRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group := createTestGroup(t, q)
	categoryRepo := NewExpenseCategoryRepo(q)
	category, err := categoryRepo.Create(ctx, group.ID, "Food")
	require.NoError(t, err)

	expense, err := expenseRepo.Create(ctx, domain.CreateExpenseParams{
		GroupID:     group.ID,
		CreatedBy:   user.ID,
		Description: "Pizza dinner",
		TotalAmount: 120.00,
		ExpenseDate: time.Now(),
		CategoryID:  &category.ID,
		SplitType:   "equal_all",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, expense.ID)
	assert.Equal(t, "Pizza dinner", expense.Description)
	assert.Equal(t, 120.00, expense.TotalAmount)
	assert.Equal(t, "equal_all", expense.SplitType)

	got, err := expenseRepo.GetByID(ctx, expense.ID)
	require.NoError(t, err)
	assert.Equal(t, expense.ID, got.ID)
}

func TestExpenseRepo_ListByGroup(t *testing.T) {
	q := newTxQueries(t)
	expenseRepo := NewExpenseRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group := createTestGroup(t, q)

	expenseRepo.Create(ctx, domain.CreateExpenseParams{
		GroupID: group.ID, CreatedBy: user.ID, Description: "Expense 1",
		TotalAmount: 50, ExpenseDate: time.Now(), SplitType: "equal_all",
	})
	expenseRepo.Create(ctx, domain.CreateExpenseParams{
		GroupID: group.ID, CreatedBy: user.ID, Description: "Expense 2",
		TotalAmount: 100, ExpenseDate: time.Now(), SplitType: "equal_all",
	})

	expenses, err := expenseRepo.ListByGroup(ctx, group.ID)
	require.NoError(t, err)
	assert.Len(t, expenses, 2)
}

func TestExpenseRepo_Update(t *testing.T) {
	q := newTxQueries(t)
	expenseRepo := NewExpenseRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group := createTestGroup(t, q)

	expense, err := expenseRepo.Create(ctx, domain.CreateExpenseParams{
		GroupID: group.ID, CreatedBy: user.ID, Description: "Original",
		TotalAmount: 100, ExpenseDate: time.Now(), SplitType: "equal_all",
	})
	require.NoError(t, err)

	newDesc := "Updated"
	newSplit := "custom"
	updated, err := expenseRepo.Update(ctx, expense.ID, domain.UpdateExpenseParams{
		Description: &newDesc,
		SplitType:   &newSplit,
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated", updated.Description)
	assert.Equal(t, "custom", updated.SplitType)
}

func TestExpenseRepo_Delete(t *testing.T) {
	q := newTxQueries(t)
	expenseRepo := NewExpenseRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group := createTestGroup(t, q)

	expense, err := expenseRepo.Create(ctx, domain.CreateExpenseParams{
		GroupID: group.ID, CreatedBy: user.ID, Description: "Delete me",
		TotalAmount: 50, ExpenseDate: time.Now(), SplitType: "equal_all",
	})
	require.NoError(t, err)

	err = expenseRepo.Delete(ctx, expense.ID)
	require.NoError(t, err)

	_, err = expenseRepo.GetByID(ctx, expense.ID)
	assert.Error(t, err)
}

func TestExpenseCategoryRepo_CreateAndList(t *testing.T) {
	q := newTxQueries(t)
	categoryRepo := NewExpenseCategoryRepo(q)
	ctx := context.Background()

	group := createTestGroup(t, q)

	cat1, err := categoryRepo.Create(ctx, group.ID, "Food")
	require.NoError(t, err)
	cat2, err := categoryRepo.Create(ctx, group.ID, "Transport")
	require.NoError(t, err)

	categories, err := categoryRepo.ListByGroup(ctx, group.ID)
	require.NoError(t, err)
	assert.Len(t, categories, 2)

	got, err := categoryRepo.GetByID(ctx, cat1.ID)
	require.NoError(t, err)
	assert.Equal(t, "Food", got.Name)

	updated, err := categoryRepo.Update(ctx, cat2.ID, group.ID, "Travel")
	require.NoError(t, err)
	assert.Equal(t, "Travel", updated.Name)
}

func TestExpenseSplitRepo_CreateMany(t *testing.T) {
	q := newTxQueries(t)
	expenseRepo := NewExpenseRepo(q)
	splitRepo := NewExpenseSplitRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group := createTestGroup(t, q)
	user2 := createTestUser(t, q)

	expense, err := expenseRepo.Create(ctx, domain.CreateExpenseParams{
		GroupID: group.ID, CreatedBy: user.ID, Description: "Split test",
		TotalAmount: 100, ExpenseDate: time.Now(), SplitType: "custom",
	})
	require.NoError(t, err)

	count, err := splitRepo.CreateMany(ctx, expense.ID, []domain.CreateExpenseSplitParams{
		{UserID: user.ID, Amount: 60},
		{UserID: user2.ID, Amount: 40},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	splits, err := splitRepo.ListByExpense(ctx, expense.ID)
	require.NoError(t, err)
	assert.Len(t, splits, 2)

	err = splitRepo.MarkAsPaid(ctx, splits[0].ID)
	require.NoError(t, err)

	gotSplit, err := splitRepo.GetByID(ctx, splits[0].ID)
	require.NoError(t, err)
	assert.True(t, gotSplit.IsPaid)
}

func TestInstallmentRepo_CreateMany(t *testing.T) {
	q := newTxQueries(t)
	expenseRepo := NewExpenseRepo(q)
	installmentRepo := NewInstallmentRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group := createTestGroup(t, q)

	expense, err := expenseRepo.Create(ctx, domain.CreateExpenseParams{
		GroupID: group.ID, CreatedBy: user.ID, Description: "Installment test",
		TotalAmount: 300, ExpenseDate: time.Now(), SplitType: "equal_all",
	})
	require.NoError(t, err)

	count, err := installmentRepo.CreateMany(ctx, expense.ID, []domain.CreateInstallmentParams{
		{InstallmentNumber: 1, Amount: 100, DueDate: time.Now().AddDate(0, 1, 0)},
		{InstallmentNumber: 2, Amount: 100, DueDate: time.Now().AddDate(0, 2, 0)},
		{InstallmentNumber: 3, Amount: 100, DueDate: time.Now().AddDate(0, 3, 0)},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	installments, err := installmentRepo.ListByExpense(ctx, expense.ID)
	require.NoError(t, err)
	assert.Len(t, installments, 3)

	err = installmentRepo.MarkAsPaid(ctx, installments[0].ID)
	require.NoError(t, err)

	got, err := installmentRepo.GetByID(ctx, installments[0].ID)
	require.NoError(t, err)
	assert.True(t, got.IsPaid)
}
