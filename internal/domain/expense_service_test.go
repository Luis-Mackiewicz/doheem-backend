package domain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockExpenseRepo struct {
	mock.Mock
}

func (m *mockExpenseRepo) GetByID(ctx context.Context, id string) (Expense, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Expense), args.Error(1)
}

func (m *mockExpenseRepo) ListByGroup(ctx context.Context, groupID string) ([]Expense, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]Expense), args.Error(1)
}

func (m *mockExpenseRepo) ListByUser(ctx context.Context, userID string) ([]Expense, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]Expense), args.Error(1)
}

func (m *mockExpenseRepo) ListByCategory(ctx context.Context, categoryID string) ([]Expense, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]Expense), args.Error(1)
}

func (m *mockExpenseRepo) Create(ctx context.Context, params CreateExpenseParams) (Expense, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(Expense), args.Error(1)
}

func (m *mockExpenseRepo) Update(ctx context.Context, id string, params UpdateExpenseParams) (Expense, error) {
	args := m.Called(ctx, id, params)
	return args.Get(0).(Expense), args.Error(1)
}

func (m *mockExpenseRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockExpenseRepo) GetTotalByGroup(ctx context.Context, groupID string) (float64, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).(float64), args.Error(1)
}

type mockExpenseSplitRepo struct {
	mock.Mock
}

func (m *mockExpenseSplitRepo) GetByID(ctx context.Context, id string) (ExpenseSplit, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ExpenseSplit), args.Error(1)
}

func (m *mockExpenseSplitRepo) ListByExpense(ctx context.Context, expenseID string) ([]ExpenseSplitWithUser, error) {
	args := m.Called(ctx, expenseID)
	return args.Get(0).([]ExpenseSplitWithUser), args.Error(1)
}

func (m *mockExpenseSplitRepo) ListByUser(ctx context.Context, userID, groupID string) ([]ExpenseSplit, error) {
	args := m.Called(ctx, userID, groupID)
	return args.Get(0).([]ExpenseSplit), args.Error(1)
}

func (m *mockExpenseSplitRepo) Create(ctx context.Context, expenseID, userID string, amount float64) (ExpenseSplit, error) {
	args := m.Called(ctx, expenseID, userID, amount)
	return args.Get(0).(ExpenseSplit), args.Error(1)
}

func (m *mockExpenseSplitRepo) CreateMany(ctx context.Context, expenseID string, splits []CreateExpenseSplitParams) (int64, error) {
	args := m.Called(ctx, expenseID, splits)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockExpenseSplitRepo) MarkAsPaid(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockExpenseSplitRepo) DeleteByExpense(ctx context.Context, expenseID string) error {
	args := m.Called(ctx, expenseID)
	return args.Error(0)
}

func (m *mockExpenseSplitRepo) GetUserBalance(ctx context.Context, userID, groupID string) (UserBalance, error) {
	args := m.Called(ctx, userID, groupID)
	return args.Get(0).(UserBalance), args.Error(1)
}

type mockInstallmentRepo struct {
	mock.Mock
}

func (m *mockInstallmentRepo) GetByID(ctx context.Context, id string) (Installment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Installment), args.Error(1)
}

func (m *mockInstallmentRepo) ListByExpense(ctx context.Context, expenseID string) ([]Installment, error) {
	args := m.Called(ctx, expenseID)
	return args.Get(0).([]Installment), args.Error(1)
}

func (m *mockInstallmentRepo) Create(ctx context.Context, params CreateInstallmentParams) (Installment, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(Installment), args.Error(1)
}

func (m *mockInstallmentRepo) CreateMany(ctx context.Context, expenseID string, installments []CreateInstallmentParams) (int64, error) {
	args := m.Called(ctx, expenseID, installments)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockInstallmentRepo) MarkAsPaid(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockInstallmentRepo) DeleteByExpense(ctx context.Context, expenseID string) error {
	args := m.Called(ctx, expenseID)
	return args.Error(0)
}

type mockCategoryRepo struct {
	mock.Mock
}

func (m *mockCategoryRepo) GetByID(ctx context.Context, id string) (ExpenseCategory, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ExpenseCategory), args.Error(1)
}

func (m *mockCategoryRepo) ListByGroup(ctx context.Context, groupID string) ([]ExpenseCategory, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]ExpenseCategory), args.Error(1)
}

func (m *mockCategoryRepo) Create(ctx context.Context, groupID, name string) (ExpenseCategory, error) {
	args := m.Called(ctx, groupID, name)
	return args.Get(0).(ExpenseCategory), args.Error(1)
}

func (m *mockCategoryRepo) Update(ctx context.Context, id, groupID, name string) (ExpenseCategory, error) {
	args := m.Called(ctx, id, groupID, name)
	return args.Get(0).(ExpenseCategory), args.Error(1)
}

func (m *mockCategoryRepo) Delete(ctx context.Context, id, groupID string) error {
	args := m.Called(ctx, id, groupID)
	return args.Error(0)
}

type mockEventBus struct {
	mock.Mock
}

func (m *mockEventBus) Publish(ctx context.Context, event DomainEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func TestExpenseService_Create_InvalidSplitTotal(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	mockSplit := new(mockExpenseSplitRepo)
	mockInstallment := new(mockInstallmentRepo)
	mockCategory := new(mockCategoryRepo)
	mockMember := new(mockGroupMemberRepo)
	svc := NewExpenseService(mockExpense, mockSplit, mockInstallment, mockCategory, mockMember, new(mockEventBus))
	ctx := context.Background()

	params := CreateExpenseWithSplitsParams{
		Expense: CreateExpenseParams{TotalAmount: 100},
		Splits:  []CreateExpenseSplitParams{{Amount: 30}, {Amount: 30}},
	}

	_, err := svc.Create(ctx, params)

	assert.ErrorIs(t, err, ErrInvalidSplitTotal)
	mockExpense.AssertNotCalled(t, "Create")
}

func TestExpenseService_Create_CategoryNotFound(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	mockSplit := new(mockExpenseSplitRepo)
	mockInstallment := new(mockInstallmentRepo)
	mockCategory := new(mockCategoryRepo)
	mockMember := new(mockGroupMemberRepo)
	svc := NewExpenseService(mockExpense, mockSplit, mockInstallment, mockCategory, mockMember, new(mockEventBus))
	ctx := context.Background()

	catID := "cat999"
	mockCategory.On("GetByID", ctx, catID).Return(ExpenseCategory{}, assert.AnError)

	params := CreateExpenseWithSplitsParams{
		Expense: CreateExpenseParams{
			TotalAmount: 100,
			CategoryID:  &catID,
		},
		Splits: []CreateExpenseSplitParams{},
	}

	_, err := svc.Create(ctx, params)

	assert.ErrorIs(t, err, ErrCategoryNotFound)
	mockExpense.AssertNotCalled(t, "Create")
}

func TestExpenseService_Create_Simple(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	mockSplit := new(mockExpenseSplitRepo)
	mockInstallment := new(mockInstallmentRepo)
	mockCategory := new(mockCategoryRepo)
	mockMember := new(mockGroupMemberRepo)
	mockBus := new(mockEventBus)
	svc := NewExpenseService(mockExpense, mockSplit, mockInstallment, mockCategory, mockMember, mockBus)
	ctx := context.Background()

	params := CreateExpenseWithSplitsParams{
		Expense: CreateExpenseParams{TotalAmount: 100},
	}
	mockExpense.On("Create", ctx, params.Expense).Return(Expense{ID: "1", TotalAmount: 100}, nil)
	mockBus.On("Publish", ctx, mock.MatchedBy(func(e DomainEvent) bool {
		return e.Type == "expense.created"
	})).Return(nil)

	expense, err := svc.Create(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, "1", expense.ID)
	mockExpense.AssertExpectations(t)
	mockBus.AssertExpectations(t)
}

func TestExpenseService_Create_Installment(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	mockSplit := new(mockExpenseSplitRepo)
	mockInstallment := new(mockInstallmentRepo)
	mockCategory := new(mockCategoryRepo)
	mockMember := new(mockGroupMemberRepo)
	mockBus := new(mockEventBus)
	svc := NewExpenseService(mockExpense, mockSplit, mockInstallment, mockCategory, mockMember, mockBus)
	ctx := context.Background()

	count := int16(3)
	dueDate := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	params := CreateExpenseWithSplitsParams{
		Expense: CreateExpenseParams{
			TotalAmount:      300,
			DueDate:          &dueDate,
			IsInstallment:    true,
			InstallmentCount: &count,
		},
	}

	mockExpense.On("Create", ctx, params.Expense).Return(Expense{ID: "1", TotalAmount: 300}, nil)
	mockInstallment.On("CreateMany", ctx, "1", mock.MatchedBy(func(inst []CreateInstallmentParams) bool {
		if len(inst) != 3 {
			return false
		}
		return inst[0].InstallmentNumber == 1 && inst[0].Amount == 100 &&
			inst[2].InstallmentNumber == 3 && inst[2].Amount == 100
	})).Return(int64(3), nil)
	mockBus.On("Publish", ctx, mock.MatchedBy(func(e DomainEvent) bool {
		return e.Type == "expense.created"
	})).Return(nil)

	expense, err := svc.Create(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, "1", expense.ID)
	mockExpense.AssertExpectations(t)
	mockInstallment.AssertExpectations(t)
	mockBus.AssertExpectations(t)
}

func TestExpenseService_GetByID_NotFound(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	svc := NewExpenseService(mockExpense, new(mockExpenseSplitRepo), new(mockInstallmentRepo), new(mockCategoryRepo), new(mockGroupMemberRepo), new(mockEventBus))
	ctx := context.Background()

	mockExpense.On("GetByID", ctx, "999").Return(Expense{}, assert.AnError)

	_, err := svc.GetByID(ctx, "999")

	assert.ErrorIs(t, err, ErrExpenseNotFound)
}

func TestExpenseService_GetTotalByGroup_Empty(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	svc := NewExpenseService(mockExpense, new(mockExpenseSplitRepo), new(mockInstallmentRepo), new(mockCategoryRepo), new(mockGroupMemberRepo), new(mockEventBus))
	ctx := context.Background()

	mockExpense.On("GetTotalByGroup", ctx, "g1").Return(0.0, nil)

	total, err := svc.GetTotalByGroup(ctx, "g1")

	assert.NoError(t, err)
	assert.Equal(t, 0.0, total)
}

func TestExpenseService_MarkSplitAsPaid(t *testing.T) {
	mockSplit := new(mockExpenseSplitRepo)
	svc := NewExpenseService(new(mockExpenseRepo), mockSplit, new(mockInstallmentRepo), new(mockCategoryRepo), new(mockGroupMemberRepo), new(mockEventBus))
	ctx := context.Background()

	mockSplit.On("MarkAsPaid", ctx, "s1").Return(nil)

	err := svc.MarkSplitAsPaid(ctx, "s1")

	assert.NoError(t, err)
}

func TestExpenseService_Create_InstallmentWithSplits(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	mockSplit := new(mockExpenseSplitRepo)
	mockInstallment := new(mockInstallmentRepo)
	mockCategory := new(mockCategoryRepo)
	mockMember := new(mockGroupMemberRepo)
	mockBus := new(mockEventBus)
	svc := NewExpenseService(mockExpense, mockSplit, mockInstallment, mockCategory, mockMember, mockBus)
	ctx := context.Background()

	count := int16(2)
	dueDate := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	splits := []CreateExpenseSplitParams{
		{UserID: "u1", Amount: 50},
		{UserID: "u2", Amount: 50},
	}
	params := CreateExpenseWithSplitsParams{
		Expense: CreateExpenseParams{
			TotalAmount:      100,
			DueDate:          &dueDate,
			IsInstallment:    true,
			InstallmentCount: &count,
		},
		Splits: splits,
	}

	mockExpense.On("Create", ctx, params.Expense).Return(Expense{ID: "1", TotalAmount: 100}, nil)
	mockInstallment.On("CreateMany", ctx, "1", mock.Anything).Return(int64(2), nil)
	mockSplit.On("CreateMany", ctx, "1", splits).Return(int64(2), nil)
	mockBus.On("Publish", ctx, mock.MatchedBy(func(e DomainEvent) bool {
		return e.Type == "expense.created"
	})).Return(nil)

	expense, err := svc.Create(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, "1", expense.ID)
	mockExpense.AssertExpectations(t)
	mockInstallment.AssertExpectations(t)
	mockSplit.AssertExpectations(t)
	mockBus.AssertExpectations(t)
}
