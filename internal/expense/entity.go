package expense

import (
	"context"
	"errors"
	"time"
)

type Expense struct {
	ID               string
	GroupID          string
	CreatedBy        string
	Description      string
	TotalAmount      float64
	ExpenseDate      time.Time
	DueDate          *time.Time
	CategoryID       *string
	SplitType        string
	IsInstallment    bool
	InstallmentCount *int16
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateExpenseParams struct {
	GroupID          string
	CreatedBy        string
	Description      string
	TotalAmount      float64
	ExpenseDate      time.Time
	DueDate          *time.Time
	CategoryID       *string
	SplitType        string
	IsInstallment    bool
	InstallmentCount *int16
}

type UpdateExpenseParams struct {
	Description *string
	TotalAmount *float64
	ExpenseDate *time.Time
	DueDate     *time.Time
	CategoryID  *string
	SplitType   *string
}

type ExpenseCategory struct {
	ID        string
	GroupID   string
	Name      string
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

type Installment struct {
	ID                string
	ExpenseID         string
	InstallmentNumber int16
	Amount            float64
	DueDate           time.Time
	IsPaid            bool
	PaidAt            *time.Time
	CreatedAt         time.Time
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
	Create(ctx context.Context, params CreateExpenseParams) (Expense, error)
	Update(ctx context.Context, id string, params UpdateExpenseParams) (Expense, error)
	Delete(ctx context.Context, id string) error
	GetTotalByGroup(ctx context.Context, groupID string) (float64, error)
}

type ExpenseCategoryRepository interface {
	GetByID(ctx context.Context, id string) (ExpenseCategory, error)
	ListByGroup(ctx context.Context, groupID string) ([]ExpenseCategory, error)
	Create(ctx context.Context, groupID, name string) (ExpenseCategory, error)
	Update(ctx context.Context, id, groupID, name string) (ExpenseCategory, error)
	Delete(ctx context.Context, id, groupID string) error
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

type InstallmentRepository interface {
	GetByID(ctx context.Context, id string) (Installment, error)
	ListByExpense(ctx context.Context, expenseID string) ([]Installment, error)
	Create(ctx context.Context, params CreateInstallmentParams) (Installment, error)
	CreateMany(ctx context.Context, expenseID string, installments []CreateInstallmentParams) (int64, error)
	MarkAsPaid(ctx context.Context, id string) error
	DeleteByExpense(ctx context.Context, expenseID string) error
}

type CreateInstallmentParams struct {
	InstallmentNumber int16
	Amount            float64
	DueDate           time.Time
}

var (
	ErrExpenseNotFound   = errors.New("expense not found")
	ErrInvalidSplitTotal = errors.New("split amounts must equal total amount")
	ErrCategoryNotFound  = errors.New("category not found")
)
