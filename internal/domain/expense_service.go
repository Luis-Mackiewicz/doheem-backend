package domain

import (
	"context"
	"errors"
)

type ExpenseService struct {
	expenseRepo      ExpenseRepository
	expenseSplitRepo ExpenseSplitRepository
	installmentRepo  InstallmentRepository
	categoryRepo     ExpenseCategoryRepository
	memberRepo       GroupMemberRepository
	eventBus         EventBus
}

func NewExpenseService(
	expenseRepo ExpenseRepository,
	expenseSplitRepo ExpenseSplitRepository,
	installmentRepo InstallmentRepository,
	categoryRepo ExpenseCategoryRepository,
	memberRepo GroupMemberRepository,
	eventBus EventBus,
) *ExpenseService {
	return &ExpenseService{
		expenseRepo:      expenseRepo,
		expenseSplitRepo: expenseSplitRepo,
		installmentRepo:  installmentRepo,
		categoryRepo:     categoryRepo,
		memberRepo:       memberRepo,
		eventBus:         eventBus,
	}
}

type CreateExpenseWithSplitsParams struct {
	Expense   CreateExpenseParams
	Splits    []CreateExpenseSplitParams
}

func (s *ExpenseService) Create(ctx context.Context, params CreateExpenseWithSplitsParams) (Expense, error) {
	var totalFromSplits float64
	for _, sp := range params.Splits {
		totalFromSplits += sp.Amount
	}
	if totalFromSplits > 0 && totalFromSplits != params.Expense.TotalAmount {
		return Expense{}, ErrInvalidSplitTotal
	}

	if params.Expense.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *params.Expense.CategoryID)
		if err != nil {
			return Expense{}, ErrCategoryNotFound
		}
	}

	if params.Expense.IsInstallment && params.Expense.InstallmentCount != nil {
		installmentAmount := params.Expense.TotalAmount / float64(*params.Expense.InstallmentCount)
		if params.Expense.DueDate == nil {
			return Expense{}, errors.New("due date is required for installment expenses")
		}
		expense, err := s.expenseRepo.Create(ctx, params.Expense)
		if err != nil {
			return Expense{}, err
		}

		installmentParams := make([]CreateInstallmentParams, *params.Expense.InstallmentCount)
		for i := int16(0); i < *params.Expense.InstallmentCount; i++ {
			dueDate := params.Expense.DueDate.AddDate(0, int(i), 0)
			installmentParams[i] = CreateInstallmentParams{
				InstallmentNumber: i + 1,
				Amount:            installmentAmount,
				DueDate:           dueDate,
			}
		}
		_, err = s.installmentRepo.CreateMany(ctx, expense.ID, installmentParams)
		if err != nil {
			return Expense{}, err
		}

		if len(params.Splits) > 0 {
			_, err = s.expenseSplitRepo.CreateMany(ctx, expense.ID, params.Splits)
			if err != nil {
				return Expense{}, err
			}
		}

		s.eventBus.Publish(ctx, DomainEvent{
			Type:     "expense.created",
			EntityID: expense.ID,
			UserID:   expense.CreatedBy,
			GroupID:  expense.GroupID,
			Payload: map[string]any{
				"total_amount": expense.TotalAmount,
				"description":  expense.Description,
			},
		})

		return expense, nil
	}

	expense, err := s.expenseRepo.Create(ctx, params.Expense)
	if err != nil {
		return Expense{}, err
	}

	if len(params.Splits) > 0 {
		_, err = s.expenseSplitRepo.CreateMany(ctx, expense.ID, params.Splits)
		if err != nil {
			return Expense{}, err
		}
	}

	s.eventBus.Publish(ctx, DomainEvent{
		Type:     "expense.created",
		EntityID: expense.ID,
		UserID:   expense.CreatedBy,
		GroupID:  expense.GroupID,
		Payload: map[string]any{
			"total_amount": expense.TotalAmount,
			"description":  expense.Description,
		},
	})

	return expense, nil
}

func (s *ExpenseService) GetByID(ctx context.Context, id string) (Expense, error) {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return Expense{}, ErrExpenseNotFound
	}
	return expense, nil
}

func (s *ExpenseService) ListByGroup(ctx context.Context, groupID string) ([]Expense, error) {
	return s.expenseRepo.ListByGroup(ctx, groupID)
}

func (s *ExpenseService) ListByUser(ctx context.Context, userID string) ([]Expense, error) {
	return s.expenseRepo.ListByUser(ctx, userID)
}

func (s *ExpenseService) ListByCategory(ctx context.Context, categoryID string) ([]Expense, error) {
	return s.expenseRepo.ListByCategory(ctx, categoryID)
}

func (s *ExpenseService) Update(ctx context.Context, id string, params UpdateExpenseParams) (Expense, error) {
	return s.expenseRepo.Update(ctx, id, params)
}

func (s *ExpenseService) Delete(ctx context.Context, id string) error {
	return s.expenseRepo.Delete(ctx, id)
}

func (s *ExpenseService) GetTotalByGroup(ctx context.Context, groupID string) (float64, error) {
	return s.expenseRepo.GetTotalByGroup(ctx, groupID)
}

func (s *ExpenseService) GetUserBalance(ctx context.Context, userID, groupID string) (UserBalance, error) {
	return s.expenseSplitRepo.GetUserBalance(ctx, userID, groupID)
}

func (s *ExpenseService) MarkSplitAsPaid(ctx context.Context, splitID string) error {
	return s.expenseSplitRepo.MarkAsPaid(ctx, splitID)
}

func (s *ExpenseService) ListSplitsByExpense(ctx context.Context, expenseID string) ([]ExpenseSplitWithUser, error) {
	return s.expenseSplitRepo.ListByExpense(ctx, expenseID)
}

func (s *ExpenseService) ListSplitsByUser(ctx context.Context, userID, groupID string) ([]ExpenseSplit, error) {
	return s.expenseSplitRepo.ListByUser(ctx, userID, groupID)
}

func (s *ExpenseService) MarkInstallmentAsPaid(ctx context.Context, installmentID string) error {
	return s.installmentRepo.MarkAsPaid(ctx, installmentID)
}

func (s *ExpenseService) ListInstallmentsByExpense(ctx context.Context, expenseID string) ([]Installment, error) {
	return s.installmentRepo.ListByExpense(ctx, expenseID)
}

func (s *ExpenseService) CreateCategory(ctx context.Context, groupID, name string) (ExpenseCategory, error) {
	return s.categoryRepo.Create(ctx, groupID, name)
}

func (s *ExpenseService) ListCategoriesByGroup(ctx context.Context, groupID string) ([]ExpenseCategory, error) {
	return s.categoryRepo.ListByGroup(ctx, groupID)
}

func (s *ExpenseService) UpdateCategory(ctx context.Context, id, groupID, name string) (ExpenseCategory, error) {
	return s.categoryRepo.Update(ctx, id, groupID, name)
}

func (s *ExpenseService) DeleteCategory(ctx context.Context, id, groupID string) error {
	return s.categoryRepo.Delete(ctx, id, groupID)
}
