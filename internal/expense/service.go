package expense

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"doheem-backend/internal/group"
	"doheem-backend/internal/notification"
)

type ExpenseService struct {
	expenseRepo      ExpenseRepository
	expenseSplitRepo ExpenseSplitRepository
	categoryRepo     ExpenseCategoryRepository
	memberRepo       group.GroupMemberRepository
	notifRepo        notification.NotificationRepository
}

func NewExpenseService(
	expenseRepo ExpenseRepository,
	expenseSplitRepo ExpenseSplitRepository,
	categoryRepo ExpenseCategoryRepository,
	memberRepo group.GroupMemberRepository,
	notifRepo notification.NotificationRepository,
) *ExpenseService {
	return &ExpenseService{
		expenseRepo:      expenseRepo,
		expenseSplitRepo: expenseSplitRepo,
		categoryRepo:     categoryRepo,
		memberRepo:       memberRepo,
		notifRepo:        notifRepo,
	}
}

func (s *ExpenseService) Create(ctx context.Context, params CreateExpenseWithSplitsParams) (Expense, error) {
	if len(params.Splits) == 0 && params.CalcParams.SplitMode != "" {
		calculated, err := s.CalculateSplits(ctx, params.CalcParams)
		if err != nil {
			return Expense{}, err
		}
		params.Splits = calculated
	}

	var totalFromSplits float64
	for _, sp := range params.Splits {
		totalFromSplits += sp.Amount
	}
	if totalFromSplits > 0 && totalFromSplits != params.Expense.Amount {
		return Expense{}, ErrInvalidSplitTotal
	}

	_, err := s.categoryRepo.GetByID(ctx, params.Expense.CategoryID)
	if err != nil {
		return Expense{}, ErrCategoryNotFound
	}

	if params.Expense.Installments > 1 {
		if params.Expense.FirstDueDate == nil {
			return Expense{}, errors.New("first_due_date is required for installment expenses")
		}

		parent := params.Expense
		parent.ParentExpenseID = nil
		parent.InstallmentIndex = nil
		parent.InstallmentTotal = nil
		expense, err := s.expenseRepo.Create(ctx, parent)
		if err != nil {
			return Expense{}, err
		}

		installmentAmount := expense.Amount / float64(expense.Installments)
		for i := int32(1); i <= expense.Installments; i++ {
			index := i
			total := expense.Installments
			childParams := params.Expense
			childParams.Amount = installmentAmount
			childParams.DueDate = expense.FirstDueDate.AddDate(0, int(i-1), 0)
			childParams.Installments = 1
			childParams.FirstDueDate = nil
			childParams.IsFixed = false
			childParams.ParentExpenseID = &expense.ID
			childParams.InstallmentIndex = &index
			childParams.InstallmentTotal = &total

			_, err := s.expenseRepo.Create(ctx, childParams)
			if err != nil {
				return Expense{}, err
			}
		}

		if len(params.Splits) > 0 {
			_, err = s.expenseSplitRepo.CreateMany(ctx, expense.ID, params.Splits)
			if err != nil {
				return Expense{}, err
			}

			for _, sp := range params.Splits {
				if sp.UserID == params.Expense.PaidBy {
					continue
				}
				title := fmt.Sprintf("Nova despesa: %s", params.Expense.Description)
				message := fmt.Sprintf("R$ %.2f (parcelada) — sua cota: R$ %.2f", params.Expense.Amount, sp.Amount)
				relatedID := &expense.ID
				s.notifRepo.Create(ctx, notification.CreateNotificationParams{
					UserID:    sp.UserID,
					Type:      "expense",
					Title:     title,
					Message:   message,
					RelatedID: relatedID,
				})
			}
		}

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

		for _, sp := range params.Splits {
			if sp.UserID == params.Expense.PaidBy {
				continue
			}
			title := fmt.Sprintf("Nova despesa: %s", params.Expense.Description)
			message := fmt.Sprintf("R$ %.2f — sua cota: R$ %.2f", params.Expense.Amount, sp.Amount)
			relatedID := &expense.ID
			s.notifRepo.Create(ctx, notification.CreateNotificationParams{
				UserID:    sp.UserID,
				Type:      "expense",
				Title:     title,
				Message:   message,
				RelatedID: relatedID,
			})
		}
	}

	return expense, nil
}

func splitEqually(total float64, userIDs []string) []CreateExpenseSplitParams {
	count := len(userIDs)
	if count == 0 {
		return nil
	}
	base := math.Floor(total*100/float64(count)) / 100
	remainder := math.Round((total-base*float64(count))*100) / 100
	splits := make([]CreateExpenseSplitParams, count)
	for i, uid := range userIDs {
		v := base
		if i == 0 {
			v = math.Round((base+remainder)*100) / 100
		}
		splits[i] = CreateExpenseSplitParams{UserID: uid, Amount: v}
	}
	return splits
}

func (s *ExpenseService) CalculateSplits(ctx context.Context, params CalculateSplitsParams) ([]CreateExpenseSplitParams, error) {
	switch params.SplitMode {
	case "equal":
		members, err := s.memberRepo.ListByGroup(ctx, params.GroupID)
		if err != nil {
			return nil, err
		}
		ids := make([]string, len(members))
		for i, m := range members {
			ids[i] = m.UserID
		}
		return splitEqually(params.Amount, ids), nil

	case "some":
		if len(params.SelectedUserIDs) < 2 {
			return nil, ErrNoSelectedMembers
		}
		return splitEqually(params.Amount, params.SelectedUserIDs), nil

	default:
		return nil, ErrInvalidSplitMode
	}
}

func (s *ExpenseService) GetByID(ctx context.Context, id string) (Expense, error) {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return Expense{}, ErrExpenseNotFound
	}
	return expense, nil
}

func (s *ExpenseService) ListByGroup(ctx context.Context, groupID string, dateFrom, dateTo *time.Time, limit, offset int32) ([]Expense, error) {
	return s.expenseRepo.ListByGroup(ctx, groupID, dateFrom, dateTo, limit, offset)
}

func (s *ExpenseService) ListByGroupFull(ctx context.Context, groupID string, dateFrom, dateTo *time.Time, search string, myUserID *string, limit, offset int32) ([]Expense, int64, error) {
	return s.expenseRepo.ListByGroupFull(ctx, groupID, dateFrom, dateTo, search, myUserID, limit, offset)
}

func (s *ExpenseService) CountByGroup(ctx context.Context, groupID string, dateFrom, dateTo *time.Time) (int, error) {
	return s.expenseRepo.CountByGroup(ctx, groupID, dateFrom, dateTo)
}

func (s *ExpenseService) ListByUser(ctx context.Context, userID string) ([]Expense, error) {
	return s.expenseRepo.ListByUser(ctx, userID)
}

func (s *ExpenseService) ListByCategory(ctx context.Context, categoryID string) ([]Expense, error) {
	return s.expenseRepo.ListByCategory(ctx, categoryID)
}

func (s *ExpenseService) ListByParent(ctx context.Context, parentID string) ([]Expense, error) {
	return s.expenseRepo.ListByParent(ctx, parentID)
}

func (s *ExpenseService) Update(ctx context.Context, id string, params UpdateExpenseParams, userID string) (Expense, error) {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return Expense{}, ErrExpenseNotFound
	}

	if expense.PaidBy != userID {
		member, err := s.memberRepo.Get(ctx, expense.GroupID, userID)
		if err != nil || !member.IsAdmin {
			return Expense{}, ErrForbidden
		}
	}

	return s.expenseRepo.Update(ctx, id, params)
}

func (s *ExpenseService) Delete(ctx context.Context, id, userID string) error {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return ErrExpenseNotFound
	}

	if expense.PaidBy != userID {
		member, err := s.memberRepo.Get(ctx, expense.GroupID, userID)
		if err != nil || !member.IsAdmin {
			return ErrForbidden
		}
	}

	hasPaid, err := s.expenseSplitRepo.HasPaidSplits(ctx, id)
	if err == nil && hasPaid {
		return ErrCannotDeleteWithPaidSplits
	}

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

func (s *ExpenseService) ListSplitsByExpenseIDs(ctx context.Context, expenseIDs []string) ([]ExpenseSplitWithUser, error) {
	return s.expenseSplitRepo.ListByExpenseIDs(ctx, expenseIDs)
}

func (s *ExpenseService) ListSplitsByUser(ctx context.Context, userID, groupID string) ([]ExpenseSplit, error) {
	return s.expenseSplitRepo.ListByUser(ctx, userID, groupID)
}

func (s *ExpenseService) CreateCategory(ctx context.Context, slug, label string) (ExpenseCategory, error) {
	return s.categoryRepo.Create(ctx, slug, label)
}

func (s *ExpenseService) ListAllCategories(ctx context.Context) ([]ExpenseCategory, error) {
	return s.categoryRepo.ListAll(ctx)
}

func (s *ExpenseService) UpdateCategory(ctx context.Context, id, slug, label string) (ExpenseCategory, error) {
	return s.categoryRepo.Update(ctx, id, slug, label)
}

func (s *ExpenseService) DeleteCategory(ctx context.Context, id string) error {
	return s.categoryRepo.Delete(ctx, id)
}
