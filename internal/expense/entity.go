package expense

import (
	"context"
	"errors"
	"time"
)

type Expense struct {
	ID               string
	GroupID          string
	Description      string
	Amount           float64
	CategoryID       string
	CompetenceDate   time.Time
	DueDate          time.Time
	PaidBy           string
	SplitMode        string
	Installments     int32
	FirstDueDate     *time.Time
	IsFixed          bool
	ParentExpenseID  *string
	InstallmentIndex *int32
	InstallmentTotal *int32
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateExpenseParams struct {
	GroupID          string
	Description      string
	Amount           float64
	CategoryID       string
	CompetenceDate   time.Time
	DueDate          time.Time
	PaidBy           string
	SplitMode        string
	Installments     int32
	FirstDueDate     *time.Time
	IsFixed          bool
	ParentExpenseID  *string
	InstallmentIndex *int32
	InstallmentTotal *int32
}

type UpdateExpenseParams struct {
	Description    *string
	Amount         *float64
	CompetenceDate *time.Time
	DueDate        *time.Time
	CategoryID     *string
	SplitMode      *string
}

type ExpenseCategory struct {
	ID        string
	Slug      string
	Label     string
	CreatedAt time.Time
}

type ExpenseSplit struct {
	ID        string
	ExpenseID string
	UserID    string
	Amount    float64
	IsPaid    bool
	PaidAt    *time.Time
	CreatedAt time.Time
}

type ExpenseSplitWithUser struct {
	ExpenseSplit
	UserName  string
	UserEmail string
}

type UserBalance struct {
	TotalOwed float64
	TotalPaid float64
}

type ExpenseRepository interface {
	GetByID(ctx context.Context, id string) (Expense, error)
	ListByGroup(ctx context.Context, groupID string) ([]Expense, error)
	ListByUser(ctx context.Context, userID string) ([]Expense, error)
	ListByCategory(ctx context.Context, categoryID string) ([]Expense, error)
	ListByParent(ctx context.Context, parentID string) ([]Expense, error)
	Create(ctx context.Context, params CreateExpenseParams) (Expense, error)
	Update(ctx context.Context, id string, params UpdateExpenseParams) (Expense, error)
	Delete(ctx context.Context, id string) error
	DeleteByParent(ctx context.Context, parentID string) error
	GetTotalByGroup(ctx context.Context, groupID string) (float64, error)
}

type ExpenseCategoryRepository interface {
	GetByID(ctx context.Context, id string) (ExpenseCategory, error)
	ListAll(ctx context.Context) ([]ExpenseCategory, error)
	Create(ctx context.Context, slug, label string) (ExpenseCategory, error)
	Update(ctx context.Context, id, slug, label string) (ExpenseCategory, error)
	Delete(ctx context.Context, id string) error
}

type ExpenseSplitRepository interface {
	GetByID(ctx context.Context, id string) (ExpenseSplit, error)
	ListByExpense(ctx context.Context, expenseID string) ([]ExpenseSplitWithUser, error)
	ListByUser(ctx context.Context, userID, groupID string) ([]ExpenseSplit, error)
	Create(ctx context.Context, expenseID, userID string, amount float64) (ExpenseSplit, error)
	CreateMany(ctx context.Context, expenseID string, splits []CreateExpenseSplitParams) (int64, error)
	MarkAsPaid(ctx context.Context, id string) error
	DeleteByExpense(ctx context.Context, expenseID string) error
	GetUserBalance(ctx context.Context, userID, groupID string) (UserBalance, error)
}

type CreateExpenseSplitParams struct {
	UserID string
	Amount float64
}

var (
	ErrExpenseNotFound          = errors.New("expense not found")
	ErrInvalidSplitTotal        = errors.New("split amounts must equal total amount")
	ErrCategoryNotFound         = errors.New("category not found")
	ErrForbidden                = errors.New("forbidden")
	ErrCannotDeleteWithPaidSplits = errors.New("cannot delete expense with paid splits")
)
