package expense

import (
	"context"
	"testing"
	"time"

	"doheem-backend/internal/group"
	"doheem-backend/internal/notification"

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

func (m *mockExpenseRepo) ListByGroup(ctx context.Context, groupID string, dateFrom, dateTo *time.Time, limit, offset int32) ([]Expense, error) {
	args := m.Called(ctx, groupID, dateFrom, dateTo, limit, offset)
	return args.Get(0).([]Expense), args.Error(1)
}

func (m *mockExpenseRepo) CountByGroup(ctx context.Context, groupID string, dateFrom, dateTo *time.Time) (int, error) {
	args := m.Called(ctx, groupID, dateFrom, dateTo)
	return args.Get(0).(int), args.Error(1)
}

func (m *mockExpenseRepo) ListByUser(ctx context.Context, userID string) ([]Expense, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]Expense), args.Error(1)
}

func (m *mockExpenseRepo) ListByCategory(ctx context.Context, categoryID string) ([]Expense, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]Expense), args.Error(1)
}

func (m *mockExpenseRepo) ListByParent(ctx context.Context, parentID string) ([]Expense, error) {
	args := m.Called(ctx, parentID)
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

func (m *mockExpenseRepo) DeleteByParent(ctx context.Context, parentID string) error {
	args := m.Called(ctx, parentID)
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

func (m *mockExpenseSplitRepo) HasPaidSplits(ctx context.Context, expenseID string) (bool, error) {
	args := m.Called(ctx, expenseID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockExpenseSplitRepo) GetUserBalance(ctx context.Context, userID, groupID string) (UserBalance, error) {
	args := m.Called(ctx, userID, groupID)
	return args.Get(0).(UserBalance), args.Error(1)
}

type mockCategoryRepo struct {
	mock.Mock
}

func (m *mockCategoryRepo) GetByID(ctx context.Context, id string) (ExpenseCategory, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ExpenseCategory), args.Error(1)
}

func (m *mockCategoryRepo) ListAll(ctx context.Context) ([]ExpenseCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]ExpenseCategory), args.Error(1)
}

func (m *mockCategoryRepo) Create(ctx context.Context, slug, label string) (ExpenseCategory, error) {
	args := m.Called(ctx, slug, label)
	return args.Get(0).(ExpenseCategory), args.Error(1)
}

func (m *mockCategoryRepo) Update(ctx context.Context, id, slug, label string) (ExpenseCategory, error) {
	args := m.Called(ctx, id, slug, label)
	return args.Get(0).(ExpenseCategory), args.Error(1)
}

func (m *mockCategoryRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockNotifRepo struct {
	mock.Mock
}

func (m *mockNotifRepo) GetByID(ctx context.Context, id string) (notification.Notification, error) {
	return notification.Notification{}, nil
}

func (m *mockNotifRepo) ListByUser(ctx context.Context, userID string, limit, offset int32) ([]notification.Notification, error) {
	return nil, nil
}

func (m *mockNotifRepo) ListUnreadByUser(ctx context.Context, userID string) ([]notification.Notification, error) {
	return nil, nil
}

func (m *mockNotifRepo) Create(ctx context.Context, params notification.CreateNotificationParams) (notification.Notification, error) {
	return notification.Notification{}, nil
}

func (m *mockNotifRepo) MarkAsRead(ctx context.Context, id, userID string) error {
	return nil
}

func (m *mockNotifRepo) MarkAllAsRead(ctx context.Context, userID string) error {
	return nil
}

func (m *mockNotifRepo) CountUnread(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func (m *mockNotifRepo) Delete(ctx context.Context, id, userID string) error {
	return nil
}

func (m *mockNotifRepo) ListByUserSearch(ctx context.Context, userID, search string, limit, offset int32) ([]notification.Notification, error) {
	return nil, nil
}

func (m *mockNotifRepo) CountByUser(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func (m *mockNotifRepo) DeleteAll(ctx context.Context, userID string) error {
	return nil
}

type mockGroupMemberRepo struct {
	mock.Mock
}

func (m *mockGroupMemberRepo) GetByID(ctx context.Context, id string) (group.GroupMember, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(group.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) Get(ctx context.Context, groupID, userID string) (group.GroupMember, error) {
	args := m.Called(ctx, groupID, userID)
	return args.Get(0).(group.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) ListByGroup(ctx context.Context, groupID string) ([]group.GroupMemberWithUser, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]group.GroupMemberWithUser), args.Error(1)
}
func (m *mockGroupMemberRepo) Create(ctx context.Context, groupID, userID string, isAdmin bool) (group.GroupMember, error) {
	args := m.Called(ctx, groupID, userID, isAdmin)
	return args.Get(0).(group.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) UpdateRole(ctx context.Context, groupID, userID string, isAdmin bool) (group.GroupMember, error) {
	args := m.Called(ctx, groupID, userID, isAdmin)
	return args.Get(0).(group.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) Remove(ctx context.Context, groupID, userID string) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}
func (m *mockGroupMemberRepo) Count(ctx context.Context, groupID string) (int64, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).(int64), args.Error(1)
}

func TestExpenseService_Create_InvalidSplitTotal(t *testing.T) {
	svc := NewExpenseService(new(mockExpenseRepo), new(mockExpenseSplitRepo), new(mockCategoryRepo), new(mockGroupMemberRepo), new(mockNotifRepo))

	ctx := context.Background()

	params := CreateExpenseWithSplitsParams{
		Expense: CreateExpenseParams{Amount: 100},
		Splits:  []CreateExpenseSplitParams{{Amount: 30}, {Amount: 30}},
	}

	_, err := svc.Create(ctx, params)

	assert.ErrorIs(t, err, ErrInvalidSplitTotal)
}

func TestExpenseService_Create_CategoryNotFound(t *testing.T) {
	mockCategory := new(mockCategoryRepo)
	svc := NewExpenseService(new(mockExpenseRepo), new(mockExpenseSplitRepo), mockCategory, new(mockGroupMemberRepo), new(mockNotifRepo))
	ctx := context.Background()

	catID := "cat999"
	mockCategory.On("GetByID", ctx, catID).Return(ExpenseCategory{}, assert.AnError)

	params := CreateExpenseWithSplitsParams{
		Expense: CreateExpenseParams{
			Amount:     100,
			CategoryID: catID,
		},
		Splits: []CreateExpenseSplitParams{},
	}

	_, err := svc.Create(ctx, params)

	assert.ErrorIs(t, err, ErrCategoryNotFound)
	mockCategory.AssertExpectations(t)
}

func TestExpenseService_Create_Simple(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	mockCategory := new(mockCategoryRepo)
	svc := NewExpenseService(mockExpense, new(mockExpenseSplitRepo), mockCategory, new(mockGroupMemberRepo), new(mockNotifRepo))
	ctx := context.Background()

	mockCategory.On("GetByID", ctx, "cat1").Return(ExpenseCategory{ID: "cat1"}, nil)

	params := CreateExpenseWithSplitsParams{
		Expense: CreateExpenseParams{
			Amount:     100,
			CategoryID: "cat1",
		},
	}
	mockExpense.On("Create", ctx, params.Expense).Return(Expense{ID: "1", Amount: 100}, nil)

	expense, err := svc.Create(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, "1", expense.ID)
	mockExpense.AssertExpectations(t)
}

func TestExpenseService_Create_Installment(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	mockCategory := new(mockCategoryRepo)
	svc := NewExpenseService(mockExpense, new(mockExpenseSplitRepo), mockCategory, new(mockGroupMemberRepo), new(mockNotifRepo))
	ctx := context.Background()

	mockCategory.On("GetByID", ctx, "cat1").Return(ExpenseCategory{ID: "cat1"}, nil)

	firstDue := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	params := CreateExpenseWithSplitsParams{
		Expense: CreateExpenseParams{
			Amount:       300,
			CategoryID:   "cat1",
			Installments: 3,
			FirstDueDate: &firstDue,
		},
	}

	parentID := "parent-1"
	mockExpense.On("Create", ctx, mock.MatchedBy(func(p CreateExpenseParams) bool {
		return p.ParentExpenseID == nil && p.Installments == 3
	})).Return(Expense{ID: parentID, Amount: 300, Installments: 3, FirstDueDate: &firstDue}, nil)

	mockExpense.On("Create", ctx, mock.MatchedBy(func(p CreateExpenseParams) bool {
		return p.ParentExpenseID != nil && *p.ParentExpenseID == parentID
	})).Return(Expense{}, nil).Times(3)

	expense, err := svc.Create(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, parentID, expense.ID)
	mockExpense.AssertExpectations(t)
}

func TestExpenseService_GetByID_NotFound(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	svc := NewExpenseService(mockExpense, new(mockExpenseSplitRepo), new(mockCategoryRepo), new(mockGroupMemberRepo), new(mockNotifRepo))
	ctx := context.Background()

	mockExpense.On("GetByID", ctx, "999").Return(Expense{}, assert.AnError)

	_, err := svc.GetByID(ctx, "999")

	assert.ErrorIs(t, err, ErrExpenseNotFound)
}

func TestExpenseService_GetTotalByGroup_Empty(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	svc := NewExpenseService(mockExpense, new(mockExpenseSplitRepo), new(mockCategoryRepo), new(mockGroupMemberRepo), new(mockNotifRepo))
	ctx := context.Background()

	mockExpense.On("GetTotalByGroup", ctx, "g1").Return(0.0, nil)

	total, err := svc.GetTotalByGroup(ctx, "g1")

	assert.NoError(t, err)
	assert.Equal(t, 0.0, total)
}

func TestExpenseService_MarkSplitAsPaid(t *testing.T) {
	mockSplit := new(mockExpenseSplitRepo)
	svc := NewExpenseService(new(mockExpenseRepo), mockSplit, new(mockCategoryRepo), new(mockGroupMemberRepo), new(mockNotifRepo))
	ctx := context.Background()

	mockSplit.On("MarkAsPaid", ctx, "s1").Return(nil)

	err := svc.MarkSplitAsPaid(ctx, "s1")

	assert.NoError(t, err)
}

func TestExpenseService_Create_InstallmentWithSplits(t *testing.T) {
	mockExpense := new(mockExpenseRepo)
	mockSplit := new(mockExpenseSplitRepo)
	mockCategory := new(mockCategoryRepo)
	svc := NewExpenseService(mockExpense, mockSplit, mockCategory, new(mockGroupMemberRepo), new(mockNotifRepo))
	ctx := context.Background()

	mockCategory.On("GetByID", ctx, "cat1").Return(ExpenseCategory{ID: "cat1"}, nil)

	firstDue := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	splits := []CreateExpenseSplitParams{
		{UserID: "u1", Amount: 50},
		{UserID: "u2", Amount: 50},
	}
	params := CreateExpenseWithSplitsParams{
		Expense: CreateExpenseParams{
			Amount:       100,
			CategoryID:   "cat1",
			Installments: 2,
			FirstDueDate: &firstDue,
		},
		Splits: splits,
	}

	parentID := "parent-1"
	mockExpense.On("Create", ctx, mock.MatchedBy(func(p CreateExpenseParams) bool {
		return p.ParentExpenseID == nil && p.Installments == 2
	})).Return(Expense{ID: parentID, Amount: 100, Installments: 2, FirstDueDate: &firstDue}, nil)

	mockExpense.On("Create", ctx, mock.MatchedBy(func(p CreateExpenseParams) bool {
		return p.ParentExpenseID != nil && *p.ParentExpenseID == parentID
	})).Return(Expense{}, nil).Times(2)

	mockSplit.On("CreateMany", ctx, parentID, splits).Return(int64(2), nil)

	expense, err := svc.Create(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, parentID, expense.ID)
	mockExpense.AssertExpectations(t)
	mockSplit.AssertExpectations(t)
}
