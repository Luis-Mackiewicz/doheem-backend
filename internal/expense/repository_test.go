package expense

import (
	"context"
	"testing"
	"time"

	"doheem-backend/internal/db"
	"doheem-backend/internal/dbtest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpenseRepo_CreateAndGetByID(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	expenseRepo := NewExpenseRepo(q)
	ctx := context.Background()

	user := dbtest.CreateTestUser(t, q)
	group := dbtest.CreateTestGroup(t, q, user)
	categoryRepo := NewExpenseCategoryRepo(q)
	category, err := categoryRepo.Create(ctx, "food", "Food")
	require.NoError(t, err)

	expense, err := expenseRepo.Create(ctx, CreateExpenseParams{
		GroupID:        group.ID,
		Description:    "Pizza dinner",
		Amount:         120.00,
		CategoryID:     category.ID,
		CompetenceDate: time.Now(),
		DueDate:        time.Now().AddDate(0, 0, 30),
		PaidBy:         user.ID,
		SplitMode:      "equal",
		Installments:   1,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, expense.ID)
	assert.Equal(t, "Pizza dinner", expense.Description)
	assert.Equal(t, 120.00, expense.Amount)
	assert.Equal(t, "equal", expense.SplitMode)

	got, err := expenseRepo.GetByID(ctx, expense.ID)
	require.NoError(t, err)
	assert.Equal(t, expense.ID, got.ID)
}

func TestExpenseRepo_ListByGroup(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	expenseRepo := NewExpenseRepo(q)
	ctx := context.Background()

	user := dbtest.CreateTestUser(t, q)
	group := dbtest.CreateTestGroup(t, q, user)
	category := createTestCategory(t, q)

	expenseRepo.Create(ctx, CreateExpenseParams{
		GroupID: group.ID, Description: "Expense 1",
		Amount: 50, CategoryID: category.ID, CompetenceDate: time.Now(),
		DueDate: time.Now().AddDate(0, 0, 30), PaidBy: user.ID, SplitMode: "equal",
		Installments: 1,
	})
	expenseRepo.Create(ctx, CreateExpenseParams{
		GroupID: group.ID, Description: "Expense 2",
		Amount: 100, CategoryID: category.ID, CompetenceDate: time.Now(),
		DueDate: time.Now().AddDate(0, 0, 30), PaidBy: user.ID, SplitMode: "equal",
		Installments: 1,
	})

	expenses, err := expenseRepo.ListByGroup(ctx, group.ID, nil, nil, 10, 0)
	require.NoError(t, err)
	assert.Len(t, expenses, 2)

	got, err := expenseRepo.ListByGroup(ctx, group.ID, datePtr(t, "2025-01-01"), datePtr(t, "2025-12-31"), 10, 0)
	require.NoError(t, err)
	assert.Len(t, got, 0)
}

func datePtr(t *testing.T, s string) *time.Time {
	t.Helper()
	d, err := time.Parse("2006-01-02", s)
	require.NoError(t, err)
	return &d
}

func TestExpenseRepo_Update(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	expenseRepo := NewExpenseRepo(q)
	ctx := context.Background()

	user := dbtest.CreateTestUser(t, q)
	group := dbtest.CreateTestGroup(t, q, user)
	category := createTestCategory(t, q)

	expense, err := expenseRepo.Create(ctx, CreateExpenseParams{
		GroupID: group.ID, Description: "Original",
		Amount: 100, CategoryID: category.ID, CompetenceDate: time.Now(),
		DueDate: time.Now().AddDate(0, 0, 30), PaidBy: user.ID, SplitMode: "equal",
		Installments: 1,
	})
	require.NoError(t, err)

	newDesc := "Updated"
	newSplit := "custom"
	updated, err := expenseRepo.Update(ctx, expense.ID, UpdateExpenseParams{
		Description: &newDesc,
		SplitMode:   &newSplit,
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated", updated.Description)
	assert.Equal(t, "custom", updated.SplitMode)
}

func TestExpenseRepo_Delete(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	expenseRepo := NewExpenseRepo(q)
	ctx := context.Background()

	user := dbtest.CreateTestUser(t, q)
	group := dbtest.CreateTestGroup(t, q, user)
	category := createTestCategory(t, q)

	expense, err := expenseRepo.Create(ctx, CreateExpenseParams{
		GroupID: group.ID, Description: "Delete me",
		Amount: 50, CategoryID: category.ID, CompetenceDate: time.Now(),
		DueDate: time.Now().AddDate(0, 0, 30), PaidBy: user.ID, SplitMode: "equal",
		Installments: 1,
	})
	require.NoError(t, err)

	err = expenseRepo.Delete(ctx, expense.ID)
	require.NoError(t, err)

	_, err = expenseRepo.GetByID(ctx, expense.ID)
	assert.Error(t, err)
}

func TestExpenseCategoryRepo_CreateAndList(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	categoryRepo := NewExpenseCategoryRepo(q)
	ctx := context.Background()

	cat1, err := categoryRepo.Create(ctx, "food", "Food")
	require.NoError(t, err)
	cat2, err := categoryRepo.Create(ctx, "transport", "Transport")
	require.NoError(t, err)

	categories, err := categoryRepo.ListAll(ctx)
	require.NoError(t, err)
	assert.Len(t, categories, 9)

	got, err := categoryRepo.GetByID(ctx, cat1.ID)
	require.NoError(t, err)
	assert.Equal(t, "food", got.Slug)
	assert.Equal(t, "Food", got.Label)

	updated, err := categoryRepo.Update(ctx, cat2.ID, "travel", "Travel")
	require.NoError(t, err)
	assert.Equal(t, "travel", updated.Slug)
	assert.Equal(t, "Travel", updated.Label)
}

func TestExpenseSplitRepo_CreateMany(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	expenseRepo := NewExpenseRepo(q)
	splitRepo := NewExpenseSplitRepo(q)
	ctx := context.Background()

	user := dbtest.CreateTestUser(t, q)
	group := dbtest.CreateTestGroup(t, q, user)
	user2 := dbtest.CreateTestUser(t, q)
	category := createTestCategory(t, q)

	expense, err := expenseRepo.Create(ctx, CreateExpenseParams{
		GroupID: group.ID, Description: "Split test",
		Amount: 100, CategoryID: category.ID, CompetenceDate: time.Now(),
		DueDate: time.Now().AddDate(0, 0, 30), PaidBy: user.ID, SplitMode: "custom",
		Installments: 1,
	})
	require.NoError(t, err)

	count, err := splitRepo.CreateMany(ctx, expense.ID, []CreateExpenseSplitParams{
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

func createTestCategory(t *testing.T, q *db.Queries) ExpenseCategory {
	ctx := context.Background()
	categoryRepo := NewExpenseCategoryRepo(q)
	categories, err := categoryRepo.ListAll(ctx)
	require.NoError(t, err)
	if len(categories) > 0 {
		return categories[0]
	}
	cat, err := categoryRepo.Create(ctx, "test", "Test")
	require.NoError(t, err)
	return cat
}
