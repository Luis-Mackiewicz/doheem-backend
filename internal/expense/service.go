package expense

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"doheem-backend/internal/group"
	"doheem-backend/internal/notification"

	"github.com/shopspring/decimal"
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
	if params.Expense.IsFixed && params.Expense.Installments > 1 {
		return Expense{}, ErrFixedWithInstallments
	}

	if len(params.Splits) == 0 && params.CalcParams.SplitMode != "" {
		calculated, err := s.CalculateSplits(ctx, params.CalcParams)
		if err != nil {
			return Expense{}, err
		}
		params.Splits = calculated
	}

	var totalFromSplits decimal.Decimal
	for _, sp := range params.Splits {
		totalFromSplits = totalFromSplits.Add(sp.Amount)
	}
	if totalFromSplits.GreaterThan(decimal.Zero) && totalFromSplits.Sub(params.Expense.Amount).Abs().GreaterThan(decimal.NewFromInt(5).Div(decimal.NewFromInt(1000))) {
		return Expense{}, ErrInvalidSplitTotal
	}

	_, err := s.categoryRepo.GetByID(ctx, params.Expense.CategoryID)
	if err != nil {
		return Expense{}, ErrCategoryNotFound
	}

	if params.Expense.Installments > 1 {
		if params.Expense.FirstDueDate == nil {
			return Expense{}, errors.New("first_due_date é obrigatória para despesas parceladas")
		}

		parent := params.Expense
		parent.ParentExpenseID = nil
		parent.InstallmentIndex = nil
		parent.InstallmentTotal = nil
		parentExpense, err := s.expenseRepo.Create(ctx, parent)
		if err != nil {
			return Expense{}, err
		}

		installmentCount := parentExpense.Installments
		installmentBase := parentExpense.Amount.Div(decimal.NewFromInt(int64(installmentCount))).Round(2)

		for i := int32(1); i <= installmentCount; i++ {
			index := i
			total := installmentCount
			childParams := params.Expense
			if i == installmentCount {
				childParams.Amount = parentExpense.Amount.Sub(installmentBase.Mul(decimal.NewFromInt(int64(installmentCount - 1))))
				childParams.Amount = childParams.Amount.Round(2)
			} else {
				childParams.Amount = installmentBase
			}
			childParams.DueDate = parentExpense.FirstDueDate.AddDate(0, int(i-1), 0)
			childParams.Installments = 1
			childParams.FirstDueDate = nil
			childParams.IsFixed = false
			childParams.ParentExpenseID = &parentExpense.ID
			childParams.InstallmentIndex = &index
			childParams.InstallmentTotal = &total

			child, err := s.expenseRepo.Create(ctx, childParams)
			if err != nil {
				return Expense{}, err
			}

			if len(params.Splits) > 0 {
				childSplits := make([]CreateExpenseSplitParams, len(params.Splits))
				for j, sp := range params.Splits {
					splitAmount := sp.Amount.Div(decimal.NewFromInt(int64(installmentCount)))
					if i == installmentCount {
						remaining := sp.Amount.Sub(splitAmount.Mul(decimal.NewFromInt(int64(installmentCount - 1))))
						splitAmount = remaining.Round(2)
					} else {
						splitAmount = splitAmount.Round(2)
					}
					childSplits[j] = CreateExpenseSplitParams{
						UserID: sp.UserID,
						Amount: splitAmount,
					}
				}

				_, err = s.expenseSplitRepo.CreateMany(ctx, child.ID, childSplits)
				if err != nil {
					return Expense{}, err
				}

				if i == 1 {
					for _, sp := range params.Splits {
						if sp.UserID == params.Expense.PaidBy {
							continue
						}
						installmentSplit := sp.Amount.Div(decimal.NewFromInt(int64(installmentCount))).Round(2)
						title := fmt.Sprintf("Nova despesa: %s", params.Expense.Description)
						message := fmt.Sprintf("R$ %s (%dx de R$ %s) — sua cota: R$ %s/mês", params.Expense.Amount.StringFixed(2), installmentCount, installmentBase.StringFixed(2), installmentSplit.StringFixed(2))
						relatedID := &child.ID
						s.notifRepo.Create(ctx, notification.CreateNotificationParams{
							UserID:    sp.UserID,
							Type:      "expense",
							Title:     title,
							Message:   message,
							RelatedID: relatedID,
						})
					}
				}
			}
		}

		return parentExpense, nil
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
			message := fmt.Sprintf("R$ %s — sua cota: R$ %s", params.Expense.Amount.StringFixed(2), sp.Amount.StringFixed(2))
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

func splitEqually(total decimal.Decimal, userIDs []string) []CreateExpenseSplitParams {
	count := len(userIDs)
	if count == 0 {
		return nil
	}
	countDec := decimal.NewFromInt(int64(count))
	base := total.Mul(decimal.NewFromInt(100)).Div(countDec).Floor().Div(decimal.NewFromInt(100))
	remainder := total.Sub(base.Mul(countDec)).Round(2)
	splits := make([]CreateExpenseSplitParams, count)
	for i, uid := range userIDs {
		v := base
		if i == 0 {
			v = base.Add(remainder).Round(2)
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

	if expense.PaidBy != userID && (expense.CreatedBy == nil || *expense.CreatedBy != userID) {
		member, err := s.memberRepo.Get(ctx, expense.GroupID, userID)
		if err != nil || !member.IsAdmin {
			return Expense{}, ErrForbidden
		}
	}

	if expense.ParentExpenseID != nil {
		return Expense{}, ErrCannotEditInstallmentChild
	}
	if expense.Installments > 1 {
		return Expense{}, ErrCannotEditInstallmentParent
	}

	hasPaid, err := s.expenseSplitRepo.HasPaidSplits(ctx, id)
	if err != nil {
		return Expense{}, err
	}
	if hasPaid {
		return Expense{}, ErrCannotEditWithPaidSplits
	}

	finalAmount := expense.Amount
	if params.Amount != nil {
		finalAmount = *params.Amount
	}
	finalMode := expense.SplitMode
	if params.SplitMode != nil {
		finalMode = *params.SplitMode
	}

	modeChanged := params.SplitMode != nil && finalMode != expense.SplitMode
	amountChanged := params.Amount != nil && !finalAmount.Equal(expense.Amount)

	if modeChanged || amountChanged {
		existingSplits, _ := s.expenseSplitRepo.ListByExpense(ctx, id)

		if err := s.expenseSplitRepo.DeleteByExpense(ctx, id); err != nil {
			return Expense{}, err
		}

		var splits []CreateExpenseSplitParams
		switch finalMode {
		case "equal":
			splits, err = s.CalculateSplits(ctx, CalculateSplitsParams{
				GroupID:   expense.GroupID,
				Amount:    finalAmount,
				SplitMode: "equal",
			})
		case "some":
			userIDs := params.SelectedUserIDs
			if len(userIDs) == 0 {
				for _, sp := range existingSplits {
					userIDs = append(userIDs, sp.UserID)
				}
			}
			if len(userIDs) < 2 {
				return Expense{}, ErrNoSelectedMembers
			}
			splits = splitEqually(finalAmount, userIDs)
		case "custom":
			customSplits := params.Splits
			if len(customSplits) == 0 && amountChanged && !modeChanged {
				for _, sp := range existingSplits {
					scaled := sp.Amount.Div(expense.Amount).Mul(finalAmount).Round(2)
					customSplits = append(customSplits, CreateExpenseSplitParams{
						UserID: sp.UserID,
						Amount: scaled,
					})
				}
				if len(customSplits) > 0 {
					var total decimal.Decimal
					for _, sp := range customSplits {
						total = total.Add(sp.Amount)
					}
					customSplits[0].Amount = customSplits[0].Amount.Add(finalAmount).Sub(total).Round(2)
				}
			}
			if len(customSplits) == 0 {
				return Expense{}, ErrInvalidSplitMode
			}
			var total decimal.Decimal
			for _, sp := range customSplits {
				total = total.Add(sp.Amount)
			}
			if total.Sub(finalAmount).Abs().GreaterThan(decimal.NewFromInt(5).Div(decimal.NewFromInt(1000))) {
				return Expense{}, ErrInvalidSplitTotal
			}
			splits = customSplits
		default:
			return Expense{}, ErrInvalidSplitMode
		}
		if err != nil {
			return Expense{}, err
		}
		if len(splits) > 0 {
			if _, err := s.expenseSplitRepo.CreateMany(ctx, id, splits); err != nil {
				return Expense{}, err
			}
		}
	}

	return s.expenseRepo.Update(ctx, id, params)
}

func (s *ExpenseService) Delete(ctx context.Context, id, userID string) error {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return ErrExpenseNotFound
	}

	if expense.PaidBy != userID && (expense.CreatedBy == nil || *expense.CreatedBy != userID) {
		member, err := s.memberRepo.Get(ctx, expense.GroupID, userID)
		if err != nil || !member.IsAdmin {
			return ErrForbidden
		}
	}

	if expense.ParentExpenseID == nil && expense.Installments > 1 {
		children, err := s.expenseRepo.ListByParent(ctx, id)
		if err != nil {
			return err
		}
		for _, child := range children {
			hasPaid, err := s.expenseSplitRepo.HasPaidSplits(ctx, child.ID)
			if err != nil {
				return err
			}
			if hasPaid {
				return ErrCannotDeleteWithPaidSplits
			}
		}
	} else {
		hasPaid, err := s.expenseSplitRepo.HasPaidSplits(ctx, id)
		if err == nil && hasPaid {
			return ErrCannotDeleteWithPaidSplits
		}
	}

	if err := s.expenseRepo.DeleteByParent(ctx, id); err != nil {
		return err
	}

	return s.expenseRepo.Delete(ctx, id)
}

func (s *ExpenseService) GetTotalByGroup(ctx context.Context, groupID string) (decimal.Decimal, error) {
	return s.expenseRepo.GetTotalByGroup(ctx, groupID)
}

func (s *ExpenseService) AutoRestoreFixedExpenses(ctx context.Context) error {
	origins, err := s.expenseRepo.ListFixedOrigins(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	for _, origin := range origins {
		count, err := s.expenseRepo.CountCloneByMonth(ctx, origin.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			slog.Warn("falha ao contar clones para despesa fixa", "id", origin.ID, "error", err)
			continue
		}
		if count > 0 {
			continue
		}

		cloneDue := sameDayInMonth(origin.DueDate, currentYear, currentMonth)
		cloneCompetence := sameDayInMonth(origin.CompetenceDate, currentYear, currentMonth)

		cloneParams := CreateExpenseParams{
			GroupID:        origin.GroupID,
			Description:    origin.Description,
			Amount:         origin.Amount,
			CategoryID:     origin.CategoryID,
			CompetenceDate: cloneCompetence,
			DueDate:        cloneDue,
			PaidBy:         origin.PaidBy,
			SplitMode:      origin.SplitMode,
			Installments:   1,
			IsFixed:        false,
			FixedSourceID:  &origin.ID,
		}

		clone, err := s.expenseRepo.Create(ctx, cloneParams)
		if err != nil {
			slog.Warn("failed to create clone for fixed expense", "id", origin.ID, "error", err)
			continue
		}

		if origin.SplitMode == "equal" || origin.SplitMode == "some" {
			members, err := s.memberRepo.ListByGroup(ctx, origin.GroupID)
			if err != nil {
				slog.Warn("failed to list members for fixed expense clone", "id", origin.ID, "error", err)
				continue
			}
			ids := make([]string, len(members))
			for i, m := range members {
				ids[i] = m.UserID
			}
			splitParams := splitEqually(origin.Amount, ids)
			s.expenseSplitRepo.CreateMany(ctx, clone.ID, splitParams)
		} else if origin.SplitMode == "custom" {
			originalSplits, err := s.expenseSplitRepo.ListByExpense(ctx, origin.ID)
			if err != nil {
				slog.Warn("failed to list splits for fixed expense", "id", origin.ID, "error", err)
				continue
			}
			cloneSplits := make([]CreateExpenseSplitParams, len(originalSplits))
			for i, sp := range originalSplits {
				cloneSplits[i] = CreateExpenseSplitParams{
					UserID: sp.UserID,
					Amount: sp.Amount,
				}
			}
			s.expenseSplitRepo.CreateMany(ctx, clone.ID, cloneSplits)
		}
	}

	return nil
}

func sameDayInMonth(source time.Time, year int, month time.Month) time.Time {
	day := source.Day()
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func (s *ExpenseService) GetUserBalance(ctx context.Context, userID, groupID string) (UserBalance, error) {
	return s.expenseSplitRepo.GetUserBalance(ctx, userID, groupID)
}

func (s *ExpenseService) GetGroupBalances(ctx context.Context, groupID string) (GroupBalance, error) {
	residents, totalDebt, err := s.expenseSplitRepo.GetGroupBalances(ctx, groupID)
	if err != nil {
		return GroupBalance{}, err
	}
	return GroupBalance{
		Residents: residents,
		TotalDebt: totalDebt,
	}, nil
}

type MarkSplitAsPaidInput struct {
	SplitID         string
	ReceiptData     *string
	ReceiptType     *string
	ReceiptFileName *string
}

func validateReceipt(receiptData, receiptType, receiptFileName *string) error {
	if receiptData == nil && receiptType == nil && receiptFileName == nil {
		return nil
	}
	if receiptData == nil || receiptType == nil || receiptFileName == nil {
		return errors.New("forneça receipt_data, receipt_type e receipt_file_name juntos")
	}
	allowedTypes := map[string]bool{
		"image/jpeg": true, "image/png": true, "image/webp": true,
		"application/pdf": true,
	}
	if !allowedTypes[*receiptType] {
		return fmt.Errorf("tipo de arquivo não permitido: %s", *receiptType)
	}
	maxSize := 5 * 1024 * 1024
	if len(*receiptData) > maxSize {
		return fmt.Errorf("arquivo muito grande: máximo %d bytes", maxSize)
	}
	return nil
}

func (s *ExpenseService) MarkSplitAsPaid(ctx context.Context, input MarkSplitAsPaidInput) error {
	if err := validateReceipt(input.ReceiptData, input.ReceiptType, input.ReceiptFileName); err != nil {
		return err
	}

	split, err := s.expenseSplitRepo.GetByID(ctx, input.SplitID)
	if err != nil {
		return ErrExpenseNotFound
	}
	if split.IsPaid {
		return ErrSplitAlreadyPaid
	}

	return s.expenseSplitRepo.MarkAsPaid(ctx, input.SplitID, input.ReceiptData, input.ReceiptType, input.ReceiptFileName)
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
